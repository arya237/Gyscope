package disk

type DiskDataSource interface {
	Read() (RawDiskState, error)
}
