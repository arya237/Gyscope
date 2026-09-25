package cpu

type LoadAverage struct {
	OneMinute      float64
	FiveMinutes    float64
	FifteenMinutes float64
}

type CPU struct {
	Usage        float64
	LogicalCores int
	LoadAverage  LoadAverage
}

