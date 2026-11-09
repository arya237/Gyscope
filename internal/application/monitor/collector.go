package monitor

import (
	"context"
	"sync"
	"time"

	"Gyscope/internal/domain/cpu"
	"Gyscope/internal/domain/disk"
	"Gyscope/internal/domain/gpu"
	"Gyscope/internal/domain/memory"
	"Gyscope/internal/domain/network"
	"Gyscope/internal/domain/process"
	"Gyscope/internal/domain/system"
)

// Reader interfaces

type CPUReader interface {
	GetState() (cpu.CPU, error)
}

type MemoryReader interface {
	GetState() (memory.Memory, error)
}

type DiskReader interface {
	GetState() (disk.Disk, error)
}

type GPUReader interface {
	GetState() ([]gpu.GPU, error)
}

type ProcessReader interface {
	GetState() ([]process.Process, error)
}

type NetworkReader interface {
	GetState() (network.Network, error)
}

// Readers contains all application use cases.

type Readers struct {
	CPU     CPUReader
	Memory  MemoryReader
	Disk    DiskReader
	GPU     GPUReader
	Process ProcessReader
	Network NetworkReader
}

// CollectionErrors stores the latest collection errors.

type CollectionErrors struct {
	CPU     error
	Memory  error
	Disk    error
	GPU     error
	Process error
	Network error
}

// Collector manages system data collection.

type Collector struct {
	readers  Readers
	interval time.Duration

	mu       sync.RWMutex
	snapshot system.Snapshot
	errors   CollectionErrors
}

// NewCollector creates a new system collector.

func NewCollector(
	readers Readers,
	interval time.Duration,
) *Collector {
	if interval <= 0 {
		interval = time.Second
	}

	return &Collector{
		readers:  readers,
		interval: interval,
	}
}

// Run starts the collection loop.
// It should be called only once for each Collector.

func (c *Collector) Run(ctx context.Context) {
	c.collect()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.collect()
		}
	}
}

// Snapshot returns the latest collected system state.

func (c *Collector) Snapshot() system.Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := c.snapshot

	snapshot.GPUs = append(
		[]gpu.GPU(nil),
		c.snapshot.GPUs...,
	)

	snapshot.Processes = append(
		[]process.Process(nil),
		c.snapshot.Processes...,
	)

	return snapshot
}

// Errors returns the latest collection errors.

func (c *Collector) Errors() CollectionErrors {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.errors
}

// collect executes all use cases sequentially.

func (c *Collector) collect() {
	next := c.Snapshot()

	var collectionErrors CollectionErrors

	// CPU
	if value, err := c.readers.CPU.GetState(); err != nil {
		collectionErrors.CPU = err
	} else {
		next.CPU = value
	}

	// Memory
	if value, err := c.readers.Memory.GetState(); err != nil {
		collectionErrors.Memory = err
	} else {
		next.Memory = value
	}

	// Disk
	if value, err := c.readers.Disk.GetState(); err != nil {
		collectionErrors.Disk = err
	} else {
		next.Disk = value
	}

	// GPU
	if value, err := c.readers.GPU.GetState(); err != nil {
		collectionErrors.GPU = err
	} else {
		next.GPUs = append(
			[]gpu.GPU(nil),
			value...,
		)
	}

	// Processes
	if value, err := c.readers.Process.GetState(); err != nil {
		collectionErrors.Process = err
	} else {
		next.Processes = append(
			[]process.Process(nil),
			value...,
		)
	}

	// Network
	if value, err := c.readers.Network.GetState(); err != nil {
		collectionErrors.Network = err
	} else {
		next.Network = value
	}

	next.CollectedAt = time.Now()

	c.mu.Lock()

	c.snapshot = next
	c.errors = collectionErrors

	c.mu.Unlock()
}