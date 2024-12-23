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

type ProcessInfo struct {
	Pid        int32  `json:"pid"`
	Name       string `json:"name"`
	Ppid       int32  `json:"ppid"`
	CreateTime int64  `json:"create_time"`
	Exe        string `json:"exe"`
	Cmd        string `json:"cmd"`
}
