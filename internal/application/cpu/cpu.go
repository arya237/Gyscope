package cpu

import(
	"Gyscope/internal/domain/cpu"
)

type CpuUseCase struct {
	source CPUDataSource
	previousCPUTimes *CPUTimes 
}


func NewCpuUseCase(source CPUDataSource) *CpuUseCase{
	return &CpuUseCase{
		source: source,
	}
}


func (u *CpuUseCase) GetState() (*cpu.CPU, error){
	current, err := u.source.Read()
	if err != nil{
		return nil, err
	}

	var usage float64

	if u.previousCPUTimes != nil{
		usage = calculateUsage(*u.previousCPUTimes, current.Times)
	}

	u.previousCPUTimes = &current.Times

	return &cpu.CPU{
		Usage: usage,
		LogicalCores: current.LogicalCores,
		LoadAverage: cpu.LoadAverage{
			OneMinute: current.LoadAverage.OneMinute,
			FiveMinutes: current.LoadAverage.FiveMinutes,
			FifteenMinutes: current.LoadAverage.FifteenMinutes,
		},
	}, nil
}


func calculateUsage(previous, current CPUTimes) float64{
	user := current.User - previous.User
	system := current.System - previous.System
	idle := current.Idle - previous.Idle
	iowait := current.IOWait - previous.IOWait
	irq := current.IRQ - previous.IRQ
	softirq := current.SoftIRQ - previous.SoftIRQ
	steal := current.Steal - previous.Steal

	total := user + system + idle + iowait + irq + softirq + steal

	if total == 0{
		return 0
	}

	idleTime := idle + iowait

	return float64(total-idleTime) / float64(total) * 100
}