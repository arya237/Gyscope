package linux

import (
	"os"
	"path/filepath"
	"strings"
)

type gpuDevice struct {
	Path     string
	VendorID string
	DeviceID string
	Vendor   string
	BusID    string
}

func discoverGPUDevices() ([]gpuDevice, error) {
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		return nil, err
	}

	devices := make([]gpuDevice, 0)

	for _, entry := range entries {
		name := entry.Name()

		if !strings.HasPrefix(name, "card") {
			continue
		}

		if strings.Contains(name, "-") {
			continue
		}

		cardPath := filepath.Join("/sys/class/drm", name)

		info, err := readPCIDeviceInfo(cardPath)
		if err != nil {
			continue
		}

		devices = append(devices, gpuDevice{
			Path:     cardPath,
			VendorID: info.VendorID,
			DeviceID: info.DeviceID,
			Vendor:   gpuVendorName(info.VendorID),
			BusID:    info.BusID,
		})
	}

	return devices, nil
}