package cpu

type RawLoadAverage struct {
    OneMinute      float64
    FiveMinutes    float64
    FifteenMinutes float64
}

type CPUTimes struct {
    User    uint64
    System  uint64
    Idle    uint64
    IOWait  uint64
    IRQ     uint64
    SoftIRQ uint64
    Steal   uint64
}

type RawCPUState struct {
    Times         CPUTimes
    LogicalCores  int
    LoadAverage   RawLoadAverage
}