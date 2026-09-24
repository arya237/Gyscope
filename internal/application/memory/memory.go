package memory

import (
	"Gyscope/internal/domain/memory"
)

type MemoryUseCase struct{
	source MemoryDataSource
}

func NewMemoryUseCase(soruce MemoryDataSource) *MemoryUseCase{
	return &MemoryUseCase{
		source: soruce,
	}
}

func (u *MemoryUseCase) GetMemoryState() (memory.Memory, error){
	current, err := u.source.Read()
	if err != nil{
		return memory.Memory{}, err
	}

	used := CalculateUsedMemory(current.Total, current.Available)
	usage := CalculateUsageMemory(used, current.Total)
	
	return memory.Memory{
		Total:     current.Total,
        Available: current.Available,
        Used:      used,
        Usage:     usage,
	}, nil
} 

func CalculateUsedMemory(totalMem uint64, availableMem uint64) uint64{
	return totalMem - availableMem
}

func CalculateUsageMemory(usedMem uint64, totalMem uint64) float64{
	return float64(usedMem) / float64(totalMem) * 100
}