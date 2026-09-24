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

	panic("not implemented")
}