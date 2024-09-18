package mutual

type CpuUtilization struct {
	Percent float64
}

type Document struct {
	Identity  string
	Timestamp int64
	Data      interface{}
}

type ListenStat struct {
	Addr string
	Port uint32
}
