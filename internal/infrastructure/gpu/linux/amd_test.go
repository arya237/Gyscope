package linux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAMDBackendReadDevice(t *testing.T) {
	root := t.TempDir()

	devicePath := filepath.Join(root, "device")

	writeTestFile(
		t,
		filepath.Join(devicePath, "product_name"),
		"AMD Radeon RX 7900 XTX\n",
	)

	writeTestFile(
		t,
		filepath.Join(devicePath, "gpu_busy_percent"),
		"42\n",
	)

	writeTestFile(
		t,
		filepath.Join(devicePath, "mem_info_vram_used"),
		"4294967296\n",
	)

	writeTestFile(
		t,
		filepath.Join(devicePath, "mem_info_vram_total"),
		"25769803776\n",
	)

	writeTestFile(
		t,
		filepath.Join(
			devicePath,
			"hwmon",
			"hwmon0",
			"temp1_input",
		),
		"55000\n",
	)

	backend := newAMDBackend()

	device := gpuDevice{
		Path:     root,
		VendorID: "0x1002",
		DeviceID: "0x744c",
		Vendor:   "AMD",
		BusID:    "0000:03:00.0",
	}

	state, err := backend.readDevice(device)
	if err != nil {
		t.Fatalf("read AMD device: %v", err)
	}

	if state.Name != "AMD Radeon RX 7900 XTX" {
		t.Errorf(
			"expected name %q, got %q",
			"AMD Radeon RX 7900 XTX",
			state.Name,
		)
	}

	if state.Vendor != "AMD" {
		t.Errorf(
			"expected vendor AMD, got %q",
			state.Vendor,
		)
	}

	if state.Usage != 42 {
		t.Errorf(
			"expected usage 42, got %.2f",
			state.Usage,
		)
	}

	if state.MemoryUsed != 4294967296 {
		t.Errorf(
			"expected memory used %d, got %d",
			uint64(4294967296),
			state.MemoryUsed,
		)
	}

	if state.MemoryTotal != 25769803776 {
		t.Errorf(
			"expected memory total %d, got %d",
			uint64(25769803776),
			state.MemoryTotal,
		)
	}

	if state.Temperature != 55 {
		t.Errorf(
			"expected temperature 55, got %.2f",
			state.Temperature,
		)
	}
}

func TestAMDBackendReadDevice_NameFallback(t *testing.T) {
	root := t.TempDir()

	devicePath := filepath.Join(root, "device")

	writeTestFile(
		t,
		filepath.Join(devicePath, "gpu_busy_percent"),
		"10\n",
	)

	writeTestFile(
		t,
		filepath.Join(devicePath, "mem_info_vram_used"),
		"100\n",
	)

	writeTestFile(
		t,
		filepath.Join(devicePath, "mem_info_vram_total"),
		"1000\n",
	)

	writeTestFile(
		t,
		filepath.Join(
			devicePath,
			"hwmon",
			"hwmon0",
			"temp1_input",
		),
		"45000\n",
	)

	backend := newAMDBackend()

	device := gpuDevice{
		Path:   root,
		Vendor: "AMD",
	}

	state, err := backend.readDevice(device)
	if err != nil {
		t.Fatalf("read AMD device: %v", err)
	}

	if state.Name != "AMD GPU" {
		t.Fatalf(
			"expected fallback name %q, got %q",
			"AMD GPU",
			state.Name,
		)
	}
}

func writeTestFile(
	t *testing.T,
	path string,
	content string,
) {
	t.Helper()

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		t.Fatalf(
			"create directory for %s: %v",
			path,
			err,
		)
	}

	if err := os.WriteFile(
		path,
		[]byte(content),
		0644,
	); err != nil {
		t.Fatalf(
			"write %s: %v",
			path,
			err,
		)
	}
}