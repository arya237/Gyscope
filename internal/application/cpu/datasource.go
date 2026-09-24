package cpu


type CPUDataSource interface {
    Read() (RawCPUState, error)
}