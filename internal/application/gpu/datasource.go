package gpu

type GPUDataSource interface {
	Read() ([]RawGPUState, error)
}
