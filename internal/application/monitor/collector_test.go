package monitor

import (
	"context"
	"testing"
	"time"

	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/gpu"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/network"
	"Gyscope/internal/domain/process"
)

type fakeCPUReader struct {
	value cpu.CPU
}

func (f fakeCPUReader) GetState() (cpu.CPU, error) {
	return f.value, nil
}

type fakeMemoryReader struct {
	value memory.Memory
}

func (f fakeMemoryReader) GetState() (memory.Memory, error) {
	return f.value, nil
}

type fakeDiskReader struct {
	value disk.Disk
}

func (f fakeDiskReader) GetState() (disk.Disk, error) {
	return f.value, nil
}

type fakeGPUReader struct {
	value []gpu.GPU
}

func (f fakeGPUReader) GetState() ([]gpu.GPU, error) {
	return f.value, nil
}

type fakeProcessReader struct {
	value []process.Process
}

func (f fakeProcessReader) GetState() ([]process.Process, error) {
	return f.value, nil
}

type fakeNetworkReader struct {
	value network.Network
}

func (f fakeNetworkReader) GetState() (network.Network, error) {
	return f.value, nil
}

func TestCollectorRunPublishesSnapshot(t *testing.T) {

	collector := NewCollector(
		Readers{
			CPU: fakeCPUReader{
				value: cpu.CPU{
					Usage: 37,
				},
			},

			Memory: fakeMemoryReader{
				value: memory.Memory{
					Usage: 62,
				},
			},

			Disk: fakeDiskReader{
				value: disk.Disk{
					Usage: 48,
				},
			},

			GPU: fakeGPUReader{
				value: []gpu.GPU{
					{
						Name: "Test GPU",
					},
				},
			},

			Process: fakeProcessReader{
				value: []process.Process{
					{
						PID: 123,
					},
				},
			},

			Network: fakeNetworkReader{
				value: network.Network{
					DownloadSpeed: 100,
				},
			},
		},
		time.Hour,
	)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	done := make(chan struct{})

	go func() {
		collector.Run(ctx)
		close(done)
	}()

	deadline := time.After(time.Second)

	for collector.Snapshot().CollectedAt.IsZero() {
		select {
		case <-deadline:
			t.Fatal("Collector did not publish snapshot")

		case <-time.After(time.Millisecond):
		}
	}

	snapshot := collector.Snapshot()

	if snapshot.CPU.Usage != 37 {
		t.Errorf("unexpected CPU usage: %v", snapshot.CPU.Usage)
	}

	if snapshot.Memory.Usage != 62 {
		t.Errorf("unexpected Memory usage: %v", snapshot.Memory.Usage)
	}

	if snapshot.Disk.Usage != 48 {
		t.Errorf("unexpected Disk usage: %v", snapshot.Disk.Usage)
	}

	if len(snapshot.GPUs) != 1 {
		t.Fatal("GPU snapshot was not collected")
	}

	if len(snapshot.Processes) != 1 {
		t.Fatal("Process snapshot was not collected")
	}

	if snapshot.Network.DownloadSpeed != 100 {
		t.Fatal("Network snapshot was not collected")
	}

	// Verify that returned slices do not expose internal storage.

	snapshot.GPUs[0].Name = "Modified GPU"
	snapshot.Processes[0].PID = 999

	freshSnapshot := collector.Snapshot()

	if freshSnapshot.GPUs[0].Name != "Test GPU" {
		t.Fatal("GPU slice was not copied")
	}

	if freshSnapshot.Processes[0].PID != 123 {
		t.Fatal("Process slice was not copied")
	}

	cancel()

	<-done
}
