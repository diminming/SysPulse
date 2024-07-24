package component

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

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
	gob.Register([]*process.Process{})
}

type Document struct {
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

	for {
		err := model.UpsertHost(map[string]interface{}{
			"host_identity": linuxId,
			"info":          info,
			"timestamp":     timestamp,
		})
		if err == nil {
			break
		}
	}

}

func saveCPUInfo(linuxId int64, infoLst []cpu.InfoStat, timestamp int64) {
	for {
		err := model.UpdateCPUInfo(map[string]interface{}{"host_identity": linuxId, "cpu_lst": infoLst, "timestamp": timestamp})
		if err == nil {
			break
		}
	}
}

func saveInterfaceInfo(linuxId int64, infoLst net.InterfaceStatList, timestamp int64) {
	for {
		err := model.UpdateInterface(map[string]interface{}{
			"host_identity": linuxId,
			"interface":     infoLst,
			"timestamp":     timestamp,
		})
		if err == nil {
			break
		}
	}
}

func CacheLinuxPort(linuxId int64, connLst []net.ConnectionStat) map[uint32]int32 {
	key := fmt.Sprintf("port_%d", linuxId)
	mapping := make(map[uint32]int32)
	for _, conn := range connLst {
		mapping[conn.Laddr.Port] = conn.Pid
	}

	entry := map[string]interface{}{}
	for key, value := range mapping {
		entry[strconv.FormatInt(int64(key), 10)] = value
	}

	model.CacheHMSet(key, entry)
	return mapping
}

func GetRemotePid(linuxId int64, port uint32) int32 {

	key := fmt.Sprintf("port_%d", linuxId)
	field := strconv.FormatUint(uint64(port), 10)
	value := model.CacheHGet(key, field)

	if value != "" {
		remotePid, _ := strconv.ParseInt(value, 10, 32)
		return int32(remotePid)
	}

	return -1
}

func isLocalIP(ip string) bool {
	return ip == "127.0.0.1"
}

func GetLinuxIdByIp(ip string) int64 {
	lst, _ := model.GetLinuxIdByIP(ip)
	if len(lst) != 1 {
		panic(fmt.Sprintf("got error when search linux id by ip: %s", ip))
	}
	return lst[0]
}

func saveConnRelation(localLinuxId int64, localPorts map[uint32]int32, connLst []net.ConnectionStat, timestamp int64) {
	for _, conn := range connLst {
		localPid := conn.Pid

		rIp := conn.Raddr.IP

		var remoteLinuxId int64
		var remotePid int32

		if isLocalIP(rIp) {
			remoteLinuxId = localLinuxId
			remotePid = localPorts[conn.Raddr.Port]
		} else {
			remoteLinuxId = GetLinuxIdByIp(rIp)
			remotePid = GetRemotePid(remoteLinuxId, conn.Raddr.Port)
		}

		if remotePid <= 0 || localPid <= 0 {
			continue
		}
		err := model.UpsertConnRelation(localLinuxId, localPid, remoteLinuxId, remotePid, timestamp)
		if err != nil {
			panic(fmt.Sprintf("Failed to execute UpsertConnRelation: %v", err))
		}
	}
}

func saveNetConnection(linuxId int64, connLst []net.ConnectionStat, timestamp int64) {
	establishedLst := make([]net.ConnectionStat, 0)
	for _, conn := range connLst {
		if conn.Status == "ESTABLISHED" {
			establishedLst = append(establishedLst, conn)
		}
	}
	lPortMapping := CacheLinuxPort(linuxId, establishedLst)
	saveConnRelation(linuxId, lPortMapping, establishedLst, timestamp)
}

func saveProc2Cache(linuxId int64, procLst []*process.Process, timestamp int64) {
	key := fmt.Sprintf("proc_%d", linuxId)

	entryLst := map[string]interface{}{}
	for _, proc := range procLst {
		procInfo := make(map[string]interface{})

		procInfo["pid"] = proc.Pid
		procInfo["name"], _ = proc.Name()
		procInfo["ppid"], _ = proc.Ppid()
		procInfo["create_time"], _ = proc.CreateTime()
		procInfo["exec"], _ = proc.Exe()

		entryLst[strconv.FormatInt(int64(proc.Pid), 10)] = common.ToString(procInfo)
	}

	entryLst["timestamp"] = strconv.FormatInt(timestamp, 10)

	model.CacheHMSet(key, entryLst)
}

func saveProc2GraphDB(linuxId int64, procLst []*process.Process, timestamp int64) []string {
	targetLst := make([]string, 0)
	for _, proc := range procLst {
		pid := proc.Pid
		name, _ := proc.Name()
		ppid, _ := proc.Ppid()
		create_time, _ := proc.CreateTime()
		exec, _ := proc.Exe()

		keyLst, _ := model.UpsertProcess(map[string]interface{}{
			"host_identity": linuxId,
			"pid":           pid,
			"info": map[string]interface{}{
				"name":        name,
				"ppid":        ppid,
				"create_time": create_time,
				"exec":        exec,
			},
			"timestamp": timestamp,
		})

		if len(keyLst) != 1 {
			log.Default().Printf("error at upsert processes, len(keyLst) = %d", len(keyLst))
			continue
		}

		targetLst = append(targetLst, keyLst[0])
	}
	return targetLst
}

func saveDeploymentRelation(linuxId int64, procKeyLst []string, timestamp int64) {
	model.UpsertDeploymentRelation(map[string]interface{}{
		"timestamp":     timestamp,
		"host_identity": linuxId,
		"procLst":       procKeyLst,
	})
}

func saveProcess(linuxId int64, procLst []*process.Process, timestamp int64) {

	saveProc2Cache(linuxId, procLst, timestamp)
	targetLst := saveProc2GraphDB(linuxId, procLst, timestamp)
	saveDeploymentRelation(linuxId, targetLst, timestamp)

}

func (hub *HubServer) saveData(doc *Document) {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Println(err)
			for i := 2; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				log.Default().Println(file, line)
			}
		}
	}()
	perfData := *doc
	values := perfData.Data
	identity := perfData.Identity
	linux := model.GetLinuxById(identity)
	switch val := values.(type) {
	case []cpu.TimesStat:
		saveCPUPerf(linux.Id, val, doc.Timestamp)
	case load.AvgStat:
		saveLoadPerf(linux.Id, &val, doc.Timestamp)
	case mem.VirtualMemoryStat:
		saveMemoryPerf(linux.Id, &val, doc.Timestamp)
	case mem.SwapMemoryStat:
		saveSwapPerf(linux.Id, &val, doc.Timestamp)
	case []disk.UsageStat:
		saveFsUsage(linux.Id, val, doc.Timestamp)
	case map[string]disk.IOCountersStat:
		saveDiskIOStat(linux.Id, val, doc.Timestamp)
	case []net.IOCountersStat:
		saveIfIOStat(linux.Id, val, doc.Timestamp)
	case host.InfoStat:
		saveHostInfo(linux.Id, val, doc.Timestamp)
	case []cpu.InfoStat:
		saveCPUInfo(linux.Id, val, doc.Timestamp)
	case net.InterfaceStatList:
		saveInterfaceInfo(linux.Id, val, doc.Timestamp)
	case []net.ConnectionStat:
		saveNetConnection(linux.Id, val, doc.Timestamp)
	case []*process.Process:
		saveProcess(linux.Id, val, doc.Timestamp)
	}
}

func (hub *HubServer) unpack() {
	for {
		buff := <-hub.buffChan

		buffer := bytes.NewBuffer(buff)
		data := new(Document)
		decoder := gob.NewDecoder(buffer)
		err := decoder.Decode(data)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		hub.saveData(data)
		// hub.mongoClient.Insert("perf_data", data)
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

func (hub *HubServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf0, _ := c.Next(-1)
	buf := make([]byte, len(buf0))
	copy(buf, buf0)
	hub.buffChan <- buf
	return gnet.None
}
