package process

type RawProcessState struct {
	PID          int
	Name         string
	CPUTime      uint64
	TotalCPUTime uint64
	Memory       uint64
}
