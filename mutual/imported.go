package mutual

type CpuUtilization struct {
	Percent float64
}

type Document struct {
	Identity  string
	Timestamp int64
	Data      interface{}
}
