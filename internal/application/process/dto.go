package process

type RawProcessState struct {
	PID          int
	Name         string
	CPUTime      uint64
	Memory       uint64
}

type RawProcessSnapshot struct {
	Processes   []RawProcessState
	TotalCPUTime uint64
}
