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

	"syspulse/tracker/linux/client"
	"syspulse/tracker/linux/common"
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
}

type Callback func()

type PerformanceData struct {
	Identity  string
	Timestamp int64
	Data      interface{}
}

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
	p := PerformanceData{
		Identity:  common.SysArgs.Identity,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	err := encoder.Encode(p)
	if err != nil {
		log.Default().Println(err)
	}
	m.client.Send(buffer.Bytes())
}

func (m *Monitor) collect_info() {
	hostInfo, _ := host.Info()
	m.Send(hostInfo)
	infoStat, _ := cpu.Info()
	m.Send(infoStat)
	ifStatLst, _ := net.Interfaces()
	m.Send(ifStatLst)
}

func (m *Monitor) stat_cpu() {

	timeStat1, _ := cpu.Times(false)
	timeStat2, _ := cpu.Times(true)

	m.Send(timeStat1)
	m.Send(timeStat2)
}

func (m *Monitor) stat_mem() {
	memStat, _ := mem.VirtualMemory()
	swapStat, _ := mem.SwapMemory()

	m.Send(memStat)
	m.Send(swapStat)
}

func (m *Monitor) stat_load() {
	loadStat, _ := load.Avg()
	m.Send(loadStat)
}

func (m *Monitor) stat_fs() {
	fsUsage := make([]disk.UsageStat, 3)
	partitionStatLst, _ := disk.Partitions(true)
	for _, partition := range partitionStatLst {
		usageStat, _ := disk.Usage(partition.Mountpoint)
		fsUsage = append(fsUsage, *usageStat)
	}
	m.Send(fsUsage)
}

func (m *Monitor) stat_disk_io() {
	ioStatMap, _ := disk.IOCounters()
	m.Send(ioStatMap)
}

func (m *Monitor) stat_if() {
	statLst, _ := net.IOCounters(true)
	m.Send(statLst)
	connLst, _ := net.Connections("all")
	m.Send(connLst)
}

func (m *Monitor) Run() {

	monitor := common.SysArgs.Monitor

	tickerStatic := time.NewTicker(time.Duration(monitor.Frequency.Static) * time.Second)
	tickerCpu := time.NewTicker(time.Duration(monitor.Frequency.Cpu) * time.Second)
	tickerLoad := time.NewTicker(time.Duration(monitor.Frequency.Load) * time.Second)
	tickerMemory := time.NewTicker(time.Duration(monitor.Frequency.Memory) * time.Second)
	tickerDiskIO := time.NewTicker(time.Duration(monitor.Frequency.DiskIO) * time.Second)
	tickerNetIf := time.NewTicker(time.Duration(monitor.Frequency.NetInterface) * time.Second)
	tickerFs := time.NewTicker(time.Duration(monitor.Frequency.FileSystem) * time.Second)
	for {
		select {
		case <-m.stopChan:
			return
		case <-tickerStatic.C:
			m.collect_info()
		case <-tickerCpu.C:
			m.stat_cpu()
		case <-tickerLoad.C:
			m.stat_load()
		case <-tickerMemory.C:
			m.stat_mem()
		case <-tickerDiskIO.C:
			m.stat_disk_io()
		case <-tickerNetIf.C:
			m.stat_if()
		case <-tickerFs.C:
			m.stat_fs()
		}
	}
}

func (m *Monitor) Stop() {
	m.stopChan <- true
	m.callback()
}
