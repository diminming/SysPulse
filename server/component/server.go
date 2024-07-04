package component

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/panjf2000/gnet/v2"

	"syspulse/common"
	"syspulse/model"
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

type PerformanceData struct {
	Identity  string
	Timestamp int64
	Data      interface{}
}

type HubServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	Addr      string
	Multicore bool
	buffChan  chan []byte
}

func NewHubServer() *HubServer {
	multicore := true
	return &HubServer{
		Addr:      common.SysArgs.Server.Hub.Addr,
		Multicore: multicore,
		buffChan:  make(chan []byte, 1000),
	}

}

func (hub *HubServer) OnBoot(eng gnet.Engine) gnet.Action {
	hub.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", hub.Multicore, hub.Addr)
	goroutine_total := 5
	for idx := 0; idx < goroutine_total; idx++ {
		go func() {
			hub.unpack()
		}()
	}
	return gnet.None
}

func saveCPUPerf(linuxId int64, lst []cpu.TimesStat, timestamp int64) {
	values := make([][]interface{}, 0, 1)
	for _, stat := range lst {
		values = append(values, []interface{}{linuxId, stat.CPU, stat.User, stat.System, stat.Idle, stat.Nice, stat.Iowait, stat.Irq, stat.Softirq, stat.Steal, stat.Guest, stat.GuestNice, timestamp})
	}
	sql := "insert into perf_cpu(`linux_id`,`cpu`,`user`,`system`,`idle`,`nice`,`iowait`,`irq`,`softirq`,`steal`,`guest`,`guestnice`,`timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	model.BulkInsert(sql, values)
}

func saveLoadPerf(linuxId int64, load *load.AvgStat, timestamp int64) {
	model.DBInsert("insert into perf_load(`linux_id`, `load1`, `load5`, `load15`, `timestamp`) value(?,?,?,?,?)", linuxId, load.Load1, load.Load5, load.Load15, timestamp)
}

func saveMemoryPerf(linuxId int64, stat *mem.VirtualMemoryStat, timestamp int64) {
	model.DBInsert("insert into perf_mem(`linux_id`, `total`,`available`,`used`,`usedpercent`,`free`,`active`,`inactive`,`wired`,`laundry`,`buffers`,`cached`,`writeback`,`dirty`,`writebacktmp`,`shared`,`slab`,`sreclaimable`,`sunreclaim`,`pagetables`,`swapcached`,`commitlimit`,`committedas`,`hightotal`,`highfree`,`lowtotal`,`lowfree`,`swaptotal`,`swapfree`,`mapped`,`vmalloctotal`,`vmallocused`,`vmallocchunk`,`hugepagestotal`,`hugepagesfree`,`hugepagesrsvd`,`hugepagessurp`,`hugepagesize`,`anonhugepages`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		linuxId, stat.Total, stat.Available, stat.Used, stat.UsedPercent, stat.Free, stat.Active, stat.Inactive, stat.Wired, stat.Laundry, stat.Buffers, stat.Cached, stat.WriteBack, stat.Dirty, stat.WriteBackTmp, stat.Shared, stat.Slab, stat.Sreclaimable, stat.Sunreclaim, stat.PageTables, stat.SwapCached, stat.CommitLimit, stat.CommittedAS, stat.HighTotal, stat.HighFree, stat.LowTotal, stat.LowFree, stat.SwapTotal, stat.SwapFree, stat.Mapped, stat.VmallocTotal, stat.VmallocUsed, stat.VmallocChunk, stat.HugePagesTotal, stat.HugePagesFree, stat.HugePagesRsvd, stat.HugePagesSurp, stat.HugePageSize, stat.AnonHugePages, timestamp)
}

func saveSwapPerf(linuxId int64, stat *mem.SwapMemoryStat, timestamp int64) {
	model.DBInsert("insert into perf_swap(`linux_id`, `total`, `used`, `free`, `usedpercent`, `sin`, `sout`, `pgin`, `pgout`, `pgfault`, `pgmajfault`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?)",
		linuxId, stat.Total, stat.Used, stat.Free, stat.UsedPercent, stat.Sin, stat.Sout, stat.PgIn, stat.PgOut, stat.PgFault, stat.PgMajFault, timestamp)
}

func saveFsUsage(linuxId int64, statLst []disk.UsageStat, timestamp int64) {
	values := make([][]interface{}, 0, 1)
	for _, stat := range statLst {
		values = append(values, []interface{}{linuxId, stat.Path, stat.Fstype, stat.Total, stat.Free, stat.Used, stat.UsedPercent, stat.InodesTotal, stat.InodesUsed, stat.InodesFree, stat.InodesUsedPercent, timestamp})
	}
	model.BulkInsert("insert into perf_fs_usage(`linux_id`, `path`, `fstype`, `total`, `free`, `used`, `usedpercent`, `inodestotal`, `inodesused`, `inodesfree`, `inodesusedpercent`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?)",
		values)
}

func saveDiskIOStat(linuxId int64, statMapping map[string]disk.IOCountersStat, timestamp int64) {
	values := make([][]interface{}, 0, 1)
	for _, stat := range statMapping {
		values = append(values, []interface{}{linuxId, stat.ReadCount, stat.MergedReadCount, stat.WriteCount, stat.MergedWriteCount, stat.ReadBytes, stat.WriteBytes, stat.ReadTime, stat.WriteTime, stat.IopsInProgress, stat.IoTime, stat.WeightedIO, stat.Name, stat.SerialNumber, stat.Label, timestamp})
	}
	model.BulkInsert("insert into perf_disk_io(`linux_id`,`readcount`,`mergedreadcount`,`writecount`,`mergedwritecount`,`readbytes`,`writebytes`,`readtime`,`writetime`,`iopsinprogress`,`iotime`,`weightedio`,`name`,`serialnumber`,`label`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		values)
}

func saveIfIOStat(linuxId int64, statLst []net.IOCountersStat, timestamp int64) {
	values := make([][]interface{}, 0, 1)
	for _, stat := range statLst {
		values = append(values, []interface{}{linuxId, stat.Name, stat.BytesSent, stat.BytesRecv, stat.PacketsSent, stat.PacketsRecv, stat.Errin, stat.Errout, stat.Dropin, stat.Dropout, stat.Fifoin, stat.Fifoout, timestamp})
	}
	model.BulkInsert("insert into perf_if_io(`linux_id`, `name`, `bytessent`, `bytesrecv`, `packetssent`, `packetsrecv`, `errin`, `errout`, `dropin`, `dropout`, `fifoin`, `fifoout`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		values)
}

func saveHostInfo(linuxId int64, info host.InfoStat, timestamp int64) {
	// model.DBInsert(
	// 	"insert into cfg_host(`linux_id`,`hostname`,`uptime`,`boottime`,`procs`,`os`,`platform`,`platformfamily`,`platformversion`,`kernelversion`,`kernelarch`,`virtualizationsystem`,`virtualizationrole`,`hostid`,`timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
	// 	linuxId, info.Hostname, info.Uptime, info.BootTime, info.Procs, info.OS, info.Platform, info.PlatformFamily, info.PlatformVersion, info.KernelVersion, info.KernelArch, info.VirtualizationSystem, info.VirtualizationRole, info.HostID, timestamp,
	// )
	model.CreateDocument("host", map[string]interface{}{
		"host_identity": linuxId,
		"info":          info,
		"timestamp":     timestamp,
	})
}

func saveCPUInfo(linuxId int64, infoLst []cpu.InfoStat, timestamp int64) {
	// values := make([][]interface{}, 0, 1)
	// for _, info := range infoLst {
	// 	values = append(values, []interface{}{linuxId, info.CPU, info.VendorID, info.Family, info.Model, info.Stepping, info.PhysicalID, info.CoreID, info.Cores, info.ModelName, info.Mhz, info.CacheSize, info.Microcode, timestamp})
	// }
	// model.BulkInsert("insert into cfg_cpu(`linux_id`, `cpu`, `vendorid`, `family`, `model`, `stepping`, `physicalid`, `coreid`, `cores`, `modelname`, `mhz`, `cachesize`, `microcode`, `timestamp`) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
	// 	values)
	model.CreateDocument("cpu", map[string]interface{}{"host_identity": linuxId, "cpu_lst": infoLst, "timestamp": timestamp})
}

func saveInterfaceInfo(linuxId int64, infoLst net.InterfaceStatList, timestamp int64) {
	// values := make([][]interface{}, 0, 1)
	// for _, info := range infoLst {
	// 	for _, addr := range info.Addrs {
	// 		values = append(values, []interface{}{linuxId, info.Index, info.Name, addr.Addr, info.HardwareAddr, info.MTU, strings.Join(info.Flags, ", "), timestamp})
	// 	}
	// 	// values = append(values, []interface{}{linuxId, info.Addrs[0].Addr, info., info.Family, info.Model, info.Stepping, info.PhysicalID, info.CoreID, info.Cores, info.ModelName, info.Mhz, info.CacheSize, info.Microcode, timestamp})
	// }
	for _, ifObj := range infoLst {
		model.CreateDocument("net_if", map[string]interface{}{"host_identity": linuxId, "if": ifObj, "timestamp": timestamp})
	}
}

func saveNetConnection(linuxId int64, connLst []net.ConnectionStat, timestamp int64) {
	for _, conn := range connLst {
		log.Default().Printf("Linux<%d> [%s:%d -> %s:%d]\n", linuxId, conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port)
	}
}

func (hub *HubServer) saveData(data *PerformanceData) {
	perfData := *data
	values := perfData.Data
	identity := perfData.Identity
	linux := model.GetLinuxById(identity)
	switch val := values.(type) {
	case []cpu.TimesStat:
		saveCPUPerf(linux.Id, val, data.Timestamp)
	case load.AvgStat:
		saveLoadPerf(linux.Id, &val, data.Timestamp)
	case mem.VirtualMemoryStat:
		saveMemoryPerf(linux.Id, &val, data.Timestamp)
	case mem.SwapMemoryStat:
		saveSwapPerf(linux.Id, &val, data.Timestamp)
	case []disk.UsageStat:
		saveFsUsage(linux.Id, val, data.Timestamp)
	case map[string]disk.IOCountersStat:
		saveDiskIOStat(linux.Id, val, data.Timestamp)
	case []net.IOCountersStat:
		saveIfIOStat(linux.Id, val, data.Timestamp)
	case host.InfoStat:
		saveHostInfo(linux.Id, val, data.Timestamp)
	case []cpu.InfoStat:
		saveCPUInfo(linux.Id, val, data.Timestamp)
	case net.InterfaceStatList:
		saveInterfaceInfo(linux.Id, val, data.Timestamp)
	case []net.ConnectionStat:
		saveNetConnection(linux.Id, val, data.Timestamp)
	}
}

func (hub *HubServer) unpack() {
	for {
		buff := <-hub.buffChan

		buffer := bytes.NewBuffer(buff)
		data := new(PerformanceData)
		decoder := gob.NewDecoder(buffer)
		err := decoder.Decode(data)
		if err != nil {
			log.Default().Println(err)
		}
		hub.saveData(data)
		// hub.mongoClient.Insert("perf_data", data)
	}
}

func (hub *HubServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf0, _ := c.Next(-1)
	buf := make([]byte, len(buf0))
	copy(buf, buf0)
	hub.buffChan <- buf
	log.Default().Printf("length of hub.buffChan: %d", len(hub.buffChan))
	return gnet.None
}
