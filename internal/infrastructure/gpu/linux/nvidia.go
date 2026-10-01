package linux

import (
	"fmt"

	applicationGPU "Gyscope/internal/application/gpu"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type nvidiaBackend struct{}

var _ GPUBackend = (*nvidiaBackend)(nil)

func newNVIDIABackend() *nvidiaBackend {
	return &nvidiaBackend{}
}

func (b *nvidiaBackend) Read(
	devices []gpuDevice,
) ([]applicationGPU.RawGPUState, error) {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf(
			"initialize NVML: %s",
			nvml.ErrorString(ret),
		)
	}

	defer nvml.Shutdown()

	result := make([]applicationGPU.RawGPUState, 0)

	for _, device := range devices {
		if device.Vendor != "NVIDIA" {
			continue
		}

		gpu, err := b.readDevice(device)
		if err != nil {
			return nil, err
		}

		result = append(result, gpu)
	}

	return result, nil
}

func (b *nvidiaBackend) readDevice(
	device gpuDevice,
) (applicationGPU.RawGPUState, error) {
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"get NVIDIA GPU count: %s",
			nvml.ErrorString(ret),
		)
	}

	for i := 0; i < count; i++ {
		nvidiaDevice, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			continue
		}

		pciInfo, ret := nvidiaDevice.GetPciInfo()
		if ret != nvml.SUCCESS {
			continue
		}

		nvmlBus := nvmlBusID(pciInfo.BusId)

		if normalizePCIAddress(nvmlBus) != normalizePCIAddress(device.BusID) {
			continue
		}

		return readNVIDIADevice(nvidiaDevice)
	}

	return applicationGPU.RawGPUState{}, fmt.Errorf(
		"NVIDIA GPU with PCI bus ID %s not found",
		device.BusID,
	)
}

func readNVIDIADevice(
	device nvml.Device,
) (applicationGPU.RawGPUState, error) {
	name, ret := device.GetName()
	if ret != nvml.SUCCESS {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"get name: %s",
			nvml.ErrorString(ret),
		)
	}

	utilization, ret := device.GetUtilizationRates()
	if ret != nvml.SUCCESS {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"get utilization: %s",
			nvml.ErrorString(ret),
		)
	}

	memory, ret := device.GetMemoryInfo()
	if ret != nvml.SUCCESS {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"get memory: %s",
			nvml.ErrorString(ret),
		)
	}

	temperature, ret := device.GetTemperature(
		nvml.TEMPERATURE_GPU,
	)
	if ret != nvml.SUCCESS {
		return applicationGPU.RawGPUState{}, fmt.Errorf(
			"get temperature: %s",
			nvml.ErrorString(ret),
		)
	}

	return applicationGPU.RawGPUState{
		Name:        name,
		Vendor:      "NVIDIA",
		Usage:       float64(utilization.Gpu),
		MemoryUsed:  memory.Used,
		MemoryTotal: memory.Total,
		Temperature: float64(temperature),
	}, nil
}

func nvmlBusID(busID [32]int8) string {
	buffer := make([]byte, 0, len(busID))

	for _, value := range busID {
		if value == 0 {
			break
		}

		buffer = append(buffer, byte(value))
	}

	return string(buffer)
}