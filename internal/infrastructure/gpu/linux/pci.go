package linux

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type pciDeviceInfo struct {
	VendorID string
	DeviceID string
	BusID    string
}

func readPCIDeviceInfo(cardPath string) (pciDeviceInfo, error) {
	devicePath := filepath.Join(cardPath, "device")

	vendor, err := readSysfsValue(filepath.Join(devicePath, "vendor"))
	if err != nil {
		return pciDeviceInfo{}, fmt.Errorf("read vendor: %w", err)
	}

	device, err := readSysfsValue(filepath.Join(devicePath, "device"))
	if err != nil {
		return pciDeviceInfo{}, fmt.Errorf("read device: %w", err)
	}

	busID, err := os.Readlink(devicePath)
	if err != nil {
		return pciDeviceInfo{}, fmt.Errorf("read device link: %w", err)
	}

	busID = filepath.Base(busID)

	return pciDeviceInfo{
		VendorID: vendor,
		DeviceID: device,
		BusID:    busID,
	}, nil
}

func readSysfsValue(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func readSysfsUint64(path string) (uint64, error) {
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


func findHWMon(devicePath string) (string, error) {
	hwmonRoot := filepath.Join(devicePath, "hwmon")

	entries, err := os.ReadDir(hwmonRoot)
	if err != nil {
		return "", fmt.Errorf(
			"read hwmon directory: %w",
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

	return "", fmt.Errorf("hwmon device not found")
}

func normalizePCIAddress(address string) string {
	address = strings.TrimSpace(address)

	parts := strings.Split(address, ":")

	switch len(parts) {
	case 2:
		// bus:function
		bus := strings.TrimLeft(parts[0], "0")
		if bus == "" {
			bus = "0"
		}

		return "0:" + bus + ":" + parts[1]

	case 3:
		// domain:bus:function
		domain := strings.TrimLeft(parts[0], "0")
		if domain == "" {
			domain = "0"
		}

		bus := strings.TrimLeft(parts[1], "0")
		if bus == "" {
			bus = "0"
		}

		return domain + ":" + bus + ":" + parts[2]

	default:
		return address
	}
}