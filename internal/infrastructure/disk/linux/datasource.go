package linux

import (
	applicationDisk "Gyscope/internal/application/disk"

	"golang.org/x/sys/unix"
)

type LinuxDiskDataSource struct{}

const (
	root = "/"
)

func NewLinuxDiskDataSource() *LinuxDiskDataSource {
	return &LinuxDiskDataSource{}
}

func (s *LinuxDiskDataSource) Read() (applicationDisk.RawDiskState, error) {
	var stat unix.Statfs_t

	if err := unix.Statfs(root, &stat); err != nil {
		return applicationDisk.RawDiskState{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)

	return applicationDisk.RawDiskState{
		Total: total,
		Free:  free,
	}, nil
}
