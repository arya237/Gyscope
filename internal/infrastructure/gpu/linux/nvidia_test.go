package linux

import "testing"

func TestNVIDIABackend(t *testing.T) {
	devices, err := discoverGPUDevices()
	if err != nil {
		t.Fatalf("discover GPU devices: %v", err)
	}

	hasNVIDIA := false

	for _, device := range devices {
		if device.Vendor == "NVIDIA" {
			hasNVIDIA = true
			break
		}
	}

	if !hasNVIDIA {
		t.Skip("no NVIDIA GPU found")
	}

	backend := newNVIDIABackend()

	gpus, err := backend.Read(devices)
	if err != nil {
		t.Fatalf("read NVIDIA GPU data: %v", err)
	}

	if len(gpus) == 0 {
		t.Fatal("expected at least one NVIDIA GPU")
	}

	for i, gpu := range gpus {
		t.Logf(
			"GPU %d: name=%s usage=%.1f%% memory=%d/%d temperature=%.1f°C",
			i,
			gpu.Name,
			gpu.Usage,
			gpu.MemoryUsed,
			gpu.MemoryTotal,
			gpu.Temperature,
		)
	}
}