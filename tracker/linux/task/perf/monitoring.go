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
	gob.Register([]*process.Process{})
	gob.Register(mutual.CpuUtilization{})
}

type Callback func()

type Monitor struct {
	client   *client.Courier
	callback Callback
	stopChan chan bool
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
	// m.client.SendPipe <- buffer.Bytes()
}

func (m *Monitor) Run() {

	monitor := common.SysArgs.Monitor

	tickerCFGHost := time.NewTicker(time.Duration(monitor.Frequency.CFGHost) * time.Second)
	tickerCFGCpu := time.NewTicker(time.Duration(monitor.Frequency.CFGCpu) * time.Second)
	tickerCFGIf := time.NewTicker(time.Duration(monitor.Frequency.CFGIf) * time.Second)

	tickerRTNetConn := time.NewTicker(time.Duration(monitor.Frequency.RTNetConn) * time.Second)
	tickerRTProc := time.NewTicker(time.Duration(monitor.Frequency.RTProc) * time.Second)

	tickerPerfCpu := time.NewTicker(time.Duration(monitor.Frequency.PerfCpu) * time.Second)
	tickerPerfLoad := time.NewTicker(time.Duration(monitor.Frequency.PerfLoad) * time.Second)
	tickerPerfMemory := time.NewTicker(time.Duration(monitor.Frequency.PerfMemory) * time.Second)
	tickerPerfNetInterface := time.NewTicker(time.Duration(monitor.Frequency.PerfNetInterface) * time.Second)
	tickerPerfDisk := time.NewTicker(time.Duration(monitor.Frequency.PerfDisk) * time.Second)
	tickerPerfFileSystem := time.NewTicker(time.Duration(monitor.Frequency.PerfFileSystem) * time.Second)
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
		case <-tickerRTNetConn.C:
			connLst, _ := net.Connections("all")
			m.Send(connLst)
		case <-tickerRTProc.C:
			procLst, _ := process.Processes()
			m.Send(procLst)
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
