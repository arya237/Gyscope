package disk

type DataSource interface {
	Read() (RawDiskState, error)
}
