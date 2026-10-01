package gpu

import domainGPU "Gyscope/internal/domain/gpu"


type GPUUsecase struct {
	datasource GPUDataSource
}

func NewGPUUsecase(datasource GPUDataSource) *GPUUsecase {
	return &GPUUsecase{
		datasource: datasource,
	}
}

func (u *GPUUsecase) GetState() ([]domainGPU.GPU, error) {
	raw, err := u.datasource.Read()
	if err != nil {
		return nil, err
	}

	gpus := make([]domainGPU.GPU, 0, len(raw))

	for _, item := range raw {
		gpus = append(gpus, domainGPU.GPU{
			Name:        item.Name,
			Usage:       item.Usage,
			MemoryUsed:  item.MemoryUsed,
			MemoryTotal: item.MemoryTotal,
			Temperature: item.Temperature,
		})
	}

	return gpus, nil
}



