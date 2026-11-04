package network


type RawNetworkInterface struct {
    Name             string
    ReceivedBytes    uint64
    TransmittedBytes uint64
}


type RawNetworkState struct {
    Interfaces []RawNetworkInterface
}