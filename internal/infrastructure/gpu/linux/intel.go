package linux

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	applicationGPU "Gyscope/internal/application/gpu"
)

type intelBackend struct {
	previousSamples map[string]intelUsageSample
}

type intelUsageSample struct {
	timestamp time.Time
	renderNs  uint64
}

type intelEngineUsage struct {
	busyNs   uint64
	capacity uint64
}

type intelFDInfo struct {
	driver  string
	pdev    string
	engines map[string]intelEngineUsage
}

var _ GPUBackend = (*intelBackend)(nil)

func newIntelBackend() *intelBackend {
	return &intelBackend{
		previousSamples: make(map[string]intelUsageSample),
	}
}

func (b *intelBackend) Read(
	devices []gpuDevice,
) ([]applicationGPU.RawGPUState, error) {
	result := make([]applicationGPU.RawGPUState, 0)

	for _, device := range devices {
		if device.Vendor != "Intel" {
			continue
		}

		gpu, err := b.readDevice(device)
		if err != nil {
			return nil, err
		}

		result = append(result, gpu)
	}

	return result, nil
}

func (b *intelBackend) readDevice(
	device gpuDevice,
) (applicationGPU.RawGPUState, error) {
	devicePath := filepath.Join(device.Path, "device")

	name := readIntelName(devicePath)

	temperature := float64(0)

	if value, err := readIntelTemperature(devicePath); err == nil {
		temperature = value
	}

	usage, err := b.readIntelUsage(device.BusID)
	if err != nil {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"read Intel GPU usage: %w",
			err,
		)
	}

	return applicationGPU.RawGPUState{
		Name:        name,
		Vendor:      "Intel",
		Usage:       usage,
		MemoryUsed:  0,
		MemoryTotal: 0,
		Temperature: temperature,
	}, nil
}

func detectIntelDriver(devicePath string) string {
	driverPath := filepath.Join(devicePath, "driver")

	target, err := os.Readlink(driverPath)
	if err != nil {
		return ""
	}

	return filepath.Base(target)
}

func readIntelName(devicePath string) string {
	productName, err := readSysfsValue(
		filepath.Join(devicePath, "product_name"),
	)
	if err == nil && productName != "" {
		return productName
	}

	return "Intel GPU"
}

func readIntelTemperature(devicePath string) (float64, error) {
	hwmonPath, err := findHWMon(devicePath)
	if err != nil {
		return 0, err
	}

	value, err := readSysfsValue(
		filepath.Join(hwmonPath, "temp1_input"),
	)
	if err != nil {
		return 0, fmt.Errorf(
			"read temp1_input: %w",
			err,
		)
	}

	millidegrees, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse temperature %q: %w",
			value,
			err,
		)
	}

	return millidegrees / 1000, nil
}

func (b *intelBackend) readIntelUsage(
	busID string,
) (float64, error) {
	now := time.Now()

	currentRenderNs, err := readIntelRenderUsage(busID)
	if err != nil {
		return 0, err
	}

	current := intelUsageSample{
		timestamp: now,
		renderNs:  currentRenderNs,
	}

	previous, exists := b.previousSamples[busID]

	b.previousSamples[busID] = current

	if !exists {
		return 0, nil
	}

	return calculateIntelUsage(previous, current), nil
}

func readIntelRenderUsage(busID string) (uint64, error) {
	entries, err := filepath.Glob("/proc/[0-9]*/fdinfo/*")
	if err != nil {
		return 0, err
	}

	var maxRenderNs uint64

	for _, path := range entries {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		info, err := parseIntelFDInfo(file, busID)
		file.Close()

		if err != nil {
			continue
		}

		engine, ok := info.engines["render"]
		if !ok {
			continue
		}

		if engine.busyNs > maxRenderNs {
			maxRenderNs = engine.busyNs
		}
	}

	return maxRenderNs, nil
}

func parseIntelFDInfo(
	reader io.Reader,
	busID string,
) (intelFDInfo, error) {
	scanner := bufio.NewScanner(reader)

	driver := ""
	pdev := ""

	result := make(map[string]intelEngineUsage)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "drm-driver":
			driver = value

		case "drm-pdev":
			pdev = value

		default:
			if strings.HasPrefix(key, "drm-engine-capacity-") {
				engineName := strings.TrimPrefix(
					key,
					"drm-engine-capacity-",
				)

				capacity, err := parseDRMUint(value)
				if err != nil {
					return intelFDInfo{}, fmt.Errorf(
						"parse %s: %w",
						key,
						err,
					)
				}

				engine := result[engineName]
				engine.capacity = capacity
				result[engineName] = engine

				continue
			}

			if strings.HasPrefix(key, "drm-engine-") {
				engineName := strings.TrimPrefix(
					key,
					"drm-engine-",
				)

				busyNs, err := parseDRMUint(value)
				if err != nil {
					return intelFDInfo{}, fmt.Errorf(
						"parse %s: %w",
						key,
						err,
					)
				}

				engine := result[engineName]
				engine.busyNs = busyNs

				if engine.capacity == 0 {
					engine.capacity = 1
				}

				result[engineName] = engine
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return intelFDInfo{}, err
	}

	if driver != "i915" && driver != "xe" {
		return intelFDInfo{}, fmt.Errorf(
			"not an Intel DRM driver: %s",
			driver,
		)
	}

	if normalizePCIAddress(pdev) != normalizePCIAddress(busID) {
		return intelFDInfo{}, fmt.Errorf(
			"DRM device mismatch: %s != %s",
			pdev,
			busID,
		)
	}

	return intelFDInfo{
		driver:  driver,
		pdev:    pdev,
		engines: result,
	}, nil
}

func parseDRMUint(value string) (uint64, error) {
	fields := strings.Fields(value)

	if len(fields) == 0 {
		return 0, fmt.Errorf("empty value")
	}

	var result uint64

	_, err := fmt.Sscan(fields[0], &result)
	if err != nil {
		return 0, fmt.Errorf(
			"parse uint %q: %w",
			fields[0],
			err,
		)
	}

	return result, nil
}

func calculateIntelUsage(
	previous intelUsageSample,
	current intelUsageSample,
) float64 {
	elapsedNs := current.timestamp.Sub(
		previous.timestamp,
	).Nanoseconds()

	if elapsedNs <= 0 {
		return 0
	}

	if current.renderNs < previous.renderNs {
		return 0
	}

	delta := current.renderNs - previous.renderNs

	usage := float64(delta) /
		float64(elapsedNs) *
		100

	if usage < 0 {
		return 0
	}

	if usage > 100 {
		return 100
	}

	return usage
}
