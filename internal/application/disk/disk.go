package disk

import (
	domainDisk "Gyscope/internal/domain/disk"
)

type DiskUseCase struct {
	dataSource DiskDataSource
}

func NewDiskUseCase(dataSource DiskDataSource) *DiskUseCase {
	return &DiskUseCase{
		dataSource: dataSource,
	}
}

func (u *DiskUseCase) GetState() (domainDisk.Disk, error) {
	raw, err := u.dataSource.Read()
	if err != nil {
		return domainDisk.Disk{}, err
	}

	used := raw.Total - raw.Free

	var usage float64

	if raw.Total > 0 {
		usage = float64(used) / float64(raw.Total) * 100
	}

	return domainDisk.Disk{
		Total: raw.Total,
		Used:  used,
		Free:  raw.Free,
		Usage: usage,
	}, nil
}
