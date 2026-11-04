package network

type NetworkDatasource interface{
	Read() (RawNetworkState, error)
}

