package process

type ProcessDataSource interface {
	Read() ([]RawProcessState, error)
}
