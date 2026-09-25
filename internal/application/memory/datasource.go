package memory

type MemoryDataSource interface{
	Read()(RawMemoryState, error)
}