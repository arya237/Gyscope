package network

import(
	"time"
	"math"
    domainNetwork "Gyscope/internal/domain/network"
)

type NetworkUseCase struct {
	source NetworkDatasource
	now func() time.Time

	previousInterfaces map[string]RawNetworkInterface
	previousSampleTime time.Time
	hasPreviousSample  bool

	downloadTotal uint64
    uploadTotal   uint64
}

func NewNetworkUseCase(source NetworkDatasource) *NetworkUseCase{
	return &NetworkUseCase{
		source: source,
		previousInterfaces: make(map[string]RawNetworkInterface),
	}
}

func (u *NetworkUseCase) GetState() (domainNetwork.Network, error) {
	current, err := u.source.Read()
	if err != nil {
		return domainNetwork.Network{}, err
	}

	sampleTime := u.now()
	currentInterfaces := make(map[string]RawNetworkInterface, len(current.Interfaces))
	for _, iface := range current.Interfaces {
		currentInterfaces[iface.Name] = iface
	}

	if !u.hasPreviousSample {
		u.previousInterfaces = currentInterfaces
		u.previousSampleTime = sampleTime
		u.hasPreviousSample = true

		return u.currentState(0, 0), nil
	}

	var downloadDelta uint64
	var uploadDelta uint64

	for name, currentInterface := range currentInterfaces {
		previousInterface, exists := u.previousInterfaces[name]
		if !exists {
			continue
		}

		if currentInterface.ReceivedBytes >= previousInterface.ReceivedBytes {
			downloadDelta = saturatingAdd(
				downloadDelta,
				currentInterface.ReceivedBytes-previousInterface.ReceivedBytes,
			)
		}

		if currentInterface.TransmittedBytes >= previousInterface.TransmittedBytes {
			uploadDelta = saturatingAdd(
				uploadDelta,
				currentInterface.TransmittedBytes-previousInterface.TransmittedBytes,
			)
		}
	}

	elapsed := sampleTime.Sub(u.previousSampleTime)

	u.downloadTotal = saturatingAdd(u.downloadTotal, downloadDelta)
	u.uploadTotal =  saturatingAdd(u.uploadTotal, uploadDelta)

	var downloadSpeed uint64
	var uploadSpeed uint64
	if elapsed > 0 {
		downloadSpeed = calculateSpeed(downloadDelta, elapsed)
		uploadSpeed = calculateSpeed(uploadDelta, elapsed)
	}

	u.previousInterfaces = currentInterfaces
	u.previousSampleTime = sampleTime

	return u.currentState(downloadSpeed, uploadSpeed), nil
}

func (u *NetworkUseCase) currentState(
	downloadSpeed uint64,
	uploadSpeed uint64,
) domainNetwork.Network {
	return domainNetwork.Network{
		DownloadSpeed: downloadSpeed,
		UploadSpeed:   uploadSpeed,
		DownloadTotal: u.downloadTotal,
		UploadTotal:   u.uploadTotal,
	}
}

func calculateSpeed(bytes uint64, elapsed time.Duration) uint64 {
	if bytes == 0 || elapsed <= 0 {
		return 0
	}

	speed := float64(bytes) / elapsed.Seconds()
	if speed >= float64(math.MaxUint64) {
		return math.MaxUint64
	}

	return uint64(speed)
}

func saturatingAdd(current, delta uint64) uint64 {
	if math.MaxUint64-current < delta {
		return math.MaxUint64
	}

	return current + delta
}