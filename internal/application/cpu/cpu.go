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


func (u *CpuUseCase) GetState() (cpu.CPU, error){

	panic("not implemented")
}