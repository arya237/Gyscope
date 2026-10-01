package linux

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	applicationGPU "Gyscope/internal/application/gpu"
)

type amdBackend struct{}

var _ GPUBackend = (*amdBackend)(nil)

func newAMDBackend() *amdBackend {
	return &amdBackend{}
}

func (b *amdBackend) Read(
	devices []gpuDevice,
) ([]applicationGPU.RawGPUState, error) {
	result := make([]applicationGPU.RawGPUState, 0)

	for _, device := range devices {
		if device.Vendor != "AMD" {
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

func (b *amdBackend) readDevice(
	device gpuDevice,
) (applicationGPU.RawGPUState, error) {
	amdgpuPath := filepath.Join(device.Path, "device")

	name := readAMDName(amdgpuPath)

	usage, err := readAMDPercentage(
		filepath.Join(amdgpuPath, "gpu_busy_percent"),
	)
	if err != nil {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"read AMD GPU usage: %w",
			err,
		)
	}

	memoryUsed, err := readAMDUint64(
		filepath.Join(amdgpuPath, "mem_info_vram_used"),
	)
	if err != nil {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"read AMD VRAM used: %w",
			err,
		)
	}

	memoryTotal, err := readAMDUint64(
		filepath.Join(amdgpuPath, "mem_info_vram_total"),
	)
	if err != nil {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"read AMD VRAM total: %w",
			err,
		)
	}

	temperature, err := readAMDTemperature(amdgpuPath)
	if err != nil {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"read AMD GPU temperature: %w",
			err,
		)
	}

	return applicationGPU.RawGPUState{
		Name:        name,
		Vendor:      "AMD",
		Usage:       usage,
		MemoryUsed:  memoryUsed,
		MemoryTotal: memoryTotal,
		Temperature: temperature,
	}, nil
}

func readAMDName(devicePath string) string {
	productName, err := readSysfsValue(
		filepath.Join(devicePath, "product_name"),
	)
	if err == nil && productName != "" {
		return productName
	}

	return "AMD GPU"
}

func readAMDPercentage(path string) (float64, error) {
	value, err := readSysfsValue(path)
	if err != nil {
		return 0, err
	}

	percentage, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse percentage %q: %w",
			value,
			err,
		)
	}

	if percentage < 0 || percentage > 100 {
		return 0, fmt.Errorf(
			"percentage out of range: %.2f",
			percentage,
		)
	}

	return percentage, nil
}

func readAMDUint64(path string) (uint64, error) {
	value, err := readSysfsValue(path)
	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse uint64 %q: %w",
			value,
			err,
		)
	}

	return result, nil
}

func readAMDTemperature(devicePath string) (float64, error) {
	hwmonPath, err := findAMDHWMon(devicePath)
	if err != nil {
		return 0, err
	}

	return readAMDTemperatureFile(
		filepath.Join(hwmonPath, "temp1_input"),
	)
}

func readAMDTemperatureFile(path string) (float64, error) {
	value, err := readSysfsValue(path)
	if err != nil {
		return 0, err
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

func findAMDHWMon(devicePath string) (string, error) {
	hwmonRoot := filepath.Join(devicePath, "hwmon")

	entries, err := os.ReadDir(hwmonRoot)
	if err != nil {
		return "", fmt.Errorf(
			"read AMD hwmon directory: %w",
			err,
		)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "hwmon") {
			return filepath.Join(hwmonRoot, entry.Name()), nil
		}
	}

	return "", fmt.Errorf(
		"AMD hwmon device not found",
	)
}