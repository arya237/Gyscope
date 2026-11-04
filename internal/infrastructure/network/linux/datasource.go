package linux

import(
	applicationNetwork "Gyscope/internal/application/network"
	"fmt"
	"os"
)

const defaultProcNetDevPath = "/proc/net/dev"

type LinuxNetworkDataSource struct{
	procNetDevPath string
}

func NewLinuxNetworkDatasource() *LinuxNetworkDataSource{
	return &LinuxNetworkDataSource{
		procNetDevPath: defaultProcNetDevPath,
	}
}

func (s *LinuxNetworkDataSource) Read() (applicationNetwork.RawNetworkState, error) {
	path := s.procNetDevPath
	if path == "" {
		path = defaultProcNetDevPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return applicationNetwork.RawNetworkState{}, fmt.Errorf("read %s: %w", path, err)
	}

	state, err := parseProcNetDev(content)
	if err != nil {
		return applicationNetwork.RawNetworkState{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return state, nil
}