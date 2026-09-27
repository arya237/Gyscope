package process

type ProcessDataSource interface {
	Read() (RawProcessSnapshot, error)
}
