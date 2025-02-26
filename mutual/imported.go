package mutual

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/syspulse/mutual/common"
)

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
	Id         string `json:"id"`
	Pid        int32  `json:"pid"`
	Name       string `json:"name"`
	Ppid       int32  `json:"ppid"`
	CreateTime int64  `json:"create_time"`
	Exe        string `json:"exe"`
	Cmd        string `json:"cmd"`
}

func (p *ProcessInfo) Hash() uint32 {
	return common.GetStringHash(fmt.Sprintf("%d|%d|%s|%s|%s|%d", p.Pid, p.Ppid, p.Name, p.Exe, p.Cmd, p.CreateTime))
}

type ProcessSnapshot struct {
	ProcessLst []ProcessInfo        `json:"process_lst"`
	ConnLst    []net.ConnectionStat `json:"connection_lst"`
}
