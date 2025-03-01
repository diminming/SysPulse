package perf

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"

	"github.com/syspulse/mutual"
	"github.com/syspulse/tracker/linux/client"
	"github.com/syspulse/tracker/linux/common"

	"github.com/shirou/gopsutil/v3/process"
)

func init() {
	gob.Register([]cpu.InfoStat{})
	gob.Register(disk.UsageStat{})
	gob.Register([]cpu.TimesStat{})
	gob.Register(load.AvgStat{})
	gob.Register(mem.VirtualMemoryStat{})
	gob.Register(mem.SwapMemoryStat{})
	gob.Register([]disk.UsageStat{})
	gob.Register(map[string]disk.IOCountersStat{})
	gob.Register([]net.IOCountersStat{})
	gob.Register(host.InfoStat{})
	gob.Register(net.InterfaceStat{})
	gob.Register([]net.ConnectionStat{})
	gob.Register(net.InterfaceStatList{})
	gob.Register([]mutual.ProcessInfo{})
	gob.Register(mutual.CpuUtilization{})
	gob.Register(mutual.ProcessSnapshot{})
}

type Callback func()

type Monitor struct {
	client   *client.Courier
	callback Callback
	stopChan chan bool
}

func GetProcessSnapshot() *mutual.ProcessSnapshot {
	procLst0, _ := process.Processes()
	procLst1 := make([]mutual.ProcessInfo, 0, len(procLst0))

	for _, proc := range procLst0 {

		ppid, _ := proc.Ppid()
		procName, _ := proc.Name()
		exe, _ := proc.Exe()
		cmd, _ := proc.Cmdline()
		createTime, _ := proc.CreateTime()

		procInfo := mutual.ProcessInfo{
			Pid:        proc.Pid,
			Name:       procName,
			Ppid:       ppid,
			CreateTime: createTime,
			Exe:        exe,
			Cmd:        cmd,
		}
		procLst1 = append(procLst1, procInfo)
	}

	connLst, _ := net.Connections("all")

	return &mutual.ProcessSnapshot{
		ProcessLst: procLst1,
		ConnLst:    connLst,
	}
}

func NewMonitor(client *client.Courier, cb Callback) (*Monitor, error) {
	monitor := Monitor{
		client:   client,
		callback: cb,
	}
	return &monitor, nil
}

func (m *Monitor) Send(data interface{}) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	p := mutual.Document{
		Identity:  common.SysArgs.Identity,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	err := encoder.Encode(p)
	if err != nil {
		log.Default().Println(err)
	}

	payload := buffer.Bytes()
	m.client.Write(payload)
}

func getDuration(item string) time.Duration {
	duration, err := time.ParseDuration(item)
	if err != nil {
		zap.L().Error("Failed to parse duration", zap.Error(err))
	}
	return duration
}

func (m *Monitor) Run() {

	monitor := common.SysArgs.Monitor
	tickerCFGHost := time.NewTicker(getDuration(monitor.Frequency.CFGHost))
	tickerCFGCpu := time.NewTicker(getDuration(monitor.Frequency.CFGCpu))
	tickerCFGIf := time.NewTicker(getDuration(monitor.Frequency.CFGIf))
	tickerRuntime := time.NewTicker(getDuration(monitor.Frequency.Runtime))
	tickerPerfCpu := time.NewTicker(getDuration(monitor.Frequency.PerfCpu))
	tickerPerfLoad := time.NewTicker(getDuration(monitor.Frequency.PerfLoad))
	tickerPerfMemory := time.NewTicker(getDuration(monitor.Frequency.PerfMemory))
	tickerPerfNetInterface := time.NewTicker(getDuration(monitor.Frequency.PerfNetInterface))
	tickerPerfDisk := time.NewTicker(getDuration(monitor.Frequency.PerfDisk))
	tickerPerfFileSystem := time.NewTicker(getDuration(monitor.Frequency.PerfFileSystem))

	for {
		select {
		case <-tickerCFGHost.C:
			hostInfo, _ := host.Info()
			m.Send(hostInfo)
		case <-tickerCFGCpu.C:
			infoStat, _ := cpu.Info()
			m.Send(infoStat)
		case <-tickerCFGIf.C:
			ifStatLst, _ := net.Interfaces()
			m.Send(ifStatLst)
		case <-tickerRuntime.C:
			snapshot := GetProcessSnapshot()
			m.Send(snapshot)

		case <-tickerPerfCpu.C:
			timeStat1, _ := cpu.Times(false)
			perc, _ := cpu.Percent(time.Second, false)

			m.Send(timeStat1)
			m.Send(mutual.CpuUtilization{
				Percent: perc[0],
			})
		case <-tickerPerfLoad.C:
			loadStat, _ := load.Avg()
			m.Send(loadStat)
		case <-tickerPerfMemory.C:
			memStat, _ := mem.VirtualMemory()
			swapStat, _ := mem.SwapMemory()

			m.Send(memStat)
			m.Send(swapStat)
		case <-tickerPerfNetInterface.C:
			statLst, _ := net.IOCounters(true)
			m.Send(statLst)
		case <-tickerPerfDisk.C:
			ioStatMap, _ := disk.IOCounters()
			m.Send(ioStatMap)
		case <-tickerPerfFileSystem.C:
			fsUsage := make([]disk.UsageStat, 3)
			partitionStatLst, _ := disk.Partitions(true)
			for _, partition := range partitionStatLst {
				usageStat, _ := disk.Usage(partition.Mountpoint)
				fsUsage = append(fsUsage, *usageStat)
			}
			m.Send(fsUsage)
		}
	}
}

func (m *Monitor) Stop() {
	m.stopChan <- true
	m.callback()
}
