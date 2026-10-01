package linux

import (
	"testing"
)

func TestLinuxGPUDataSource(t *testing.T) {
	dataSource := NewLinuxGPUDataSource()

	gpus, err := dataSource.Read()
	if err != nil {
		t.Fatalf("read GPU data: %v", err)
	}

	for i, gpu := range gpus {
		t.Logf(
			"GPU %d: name=%s vendor=%s usage=%.1f%% memory=%d/%d temperature=%.1f°C",
			i,
			gpu.Name,
			gpu.Vendor,
			gpu.Usage,
			gpu.MemoryUsed,
			gpu.MemoryTotal,
			gpu.Temperature,
		)
	}
}

