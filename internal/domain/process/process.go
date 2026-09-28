package process

type Process struct {
	PID         int
	Name        string
	CPUUsage    float64
	Memory      uint64
	MemoryUsage float64
}
