package component

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	gonet "net"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"

	ants "github.com/panjf2000/ants/v2"
	gnet "github.com/panjf2000/gnet/v2"

	"github.com/syspulse/common"
	"github.com/syspulse/model"

	"github.com/syspulse/mutual"
)

var pool4PerfData, pool4Alarm *ants.PoolWithFunc

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

type HubServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	Addr      string
	Multicore bool
	batch     uint32
	coreQueue chan []byte
}

func (hub *HubServer) OnBoot(eng gnet.Engine) gnet.Action {
	hub.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", hub.Multicore, hub.Addr)
	go hub.handler()

	return gnet.None
}

func (hub *HubServer) handler() {
	for {
		data, ok := <-hub.coreQueue
		if ok {
			doc := new(mutual.Document)
			decoder := gob.NewDecoder(bytes.NewReader(data))
			err := decoder.Decode(doc)
			if err != nil {
				log.Default().Println(err)
			}

			pool4PerfData.Invoke(doc)
			pool4Alarm.Invoke(doc)
		}

	}
}

func (hub *HubServer) OnTraffic(c gnet.Conn) gnet.Action {
	for i := 0; i < 10; i++ {
		data, err := c.Peek(5)
		if err != nil {
			if errors.Is(err, io.ErrShortBuffer) {
				break
			} else {
				zap.L().Error("can't read data from conn: %v", zap.Error(err))
				return gnet.Close
			}
		}

		if data[0] != 'S' {
			zap.L().Warn("got a packet does not conform to the rules...")
			return gnet.Close
		}

		dataLength := int(binary.LittleEndian.Uint32(data[1:]))
		msgLength := dataLength + 5
		data, err = c.Peek(msgLength)
		if err != nil {
			if errors.Is(err, io.ErrShortBuffer) {
				break
			} else {
				zap.L().Error("can't read valid data from conn: %v", zap.Error(err))
				return gnet.Close
			}
		}

		if len(data) < msgLength {
			zap.L().Error("unexpected short data: ", zap.Int("expected", msgLength), zap.Int("actual", len(data)))
			return gnet.None
		}

		body := make([]byte, msgLength)
		copy(body, data[5:msgLength])
		c.Discard(msgLength)

		hub.coreQueue <- body

	}

	return gnet.None
}

func NewHubServer() *HubServer {

	pool0, err := ants.NewPoolWithFunc(1000, dispatch)

	if err != nil {
		log.Default().Fatalf("error create routine pool: %v", err)
	}

	pool1, err := ants.NewPoolWithFunc(1000, check)

	if err != nil {
		log.Default().Fatalf("error create routine pool: %v", err)
	}

	pool4PerfData = pool0
	pool4Alarm = pool1

	return &HubServer{
		Addr:      common.SysArgs.Server.Hub.Addr,
		Multicore: true,
		batch:     common.SysArgs.Server.Hub.BatchSize,
		coreQueue: make(chan []byte, common.SysArgs.Server.Hub.QueueSize),
	}

}

func check(arg any) {
	doc := arg.(*mutual.Document)
	data := doc.Data
	switch val := data.(type) {
	case mutual.CpuUtilization:
		obj := model.PerfData{
			CPU: model.CpuPerfData{
				CpuUtil: val.Percent,
			},
			Subject: doc.Identity,
		}
		// zap.L().Debug("got cpu perc data.", zap.Float64("percent", val.Percent))
		TriggerCheck(doc.Identity, obj, model.DataType_CpuPerformence, doc.Timestamp)
	case load.AvgStat:
		obj := model.PerfData{
			Load: model.LoadPerfData{
				Load1:  val.Load1,
				Load5:  val.Load5,
				Load15: val.Load15,
			},
			Subject: doc.Identity,
		}
		TriggerCheck(doc.Identity, obj, model.DataType_LoadPerformence, doc.Timestamp)
	case mem.VirtualMemoryStat:
		obj := model.PerfData{
			Memory: model.MemoryPerfData{
				Total:          val.Total,
				Free:           val.Free,
				Active:         val.Active,
				Inactive:       val.Inactive,
				Wired:          val.Wired,
				Laundry:        val.Laundry,
				Buffers:        val.Buffers,
				Cached:         val.Cached,
				WriteBack:      val.WriteBack,
				Dirty:          val.Dirty,
				WriteBackTmp:   val.WriteBackTmp,
				Shared:         val.Shared,
				Slab:           val.Slab,
				Sreclaimable:   val.Sreclaimable,
				Sunreclaim:     val.Sunreclaim,
				PageTables:     val.PageTables,
				SwapCached:     val.SwapCached,
				CommitLimit:    val.CommitLimit,
				CommittedAS:    val.CommittedAS,
				HighTotal:      val.HighTotal,
				HighFree:       val.HighFree,
				LowTotal:       val.LowTotal,
				LowFree:        val.LowFree,
				SwapTotal:      val.SwapTotal,
				SwapFree:       val.SwapFree,
				Mapped:         val.Mapped,
				VmallocTotal:   val.VmallocTotal,
				VmallocUsed:    val.VmallocUsed,
				VmallocChunk:   val.VmallocChunk,
				HugePagesTotal: val.HugePagesTotal,
				HugePagesFree:  val.HugePagesFree,
				HugePagesRsvd:  val.HugePagesRsvd,
				HugePagesSurp:  val.HugePagesSurp,
				HugePageSize:   val.HugePageSize,
				AnonHugePages:  val.AnonHugePages,
			},
			Subject: doc.Identity,
		}
		TriggerCheck(doc.Identity, obj, model.DataType_MemoryPerformence, doc.Timestamp)
	case mem.SwapMemoryStat:
		obj := model.PerfData{
			Swap: model.SwapPerfData{
				Total:       val.Total,
				Used:        val.Used,
				Free:        val.Free,
				UsedPercent: val.UsedPercent,
				Sin:         val.Sin,
				Sout:        val.Sout,
				PgIn:        val.PgIn,
				PgOut:       val.PgOut,
				PgFault:     val.PgFault,
				PgMajFault:  val.PgMajFault,
			},
			Subject: doc.Identity,
		}
		TriggerCheck(doc.Identity, obj, model.DataType_SwapPerformence, doc.Timestamp)
	case []disk.UsageStat:
		for _, item := range val {
			obj := model.PerfData{
				Disk: model.DiskPerfData{
					Path:              item.Path,
					Fstype:            item.Fstype,
					Total:             item.Total,
					Free:              item.Free,
					Used:              item.Used,
					UsedPercent:       item.UsedPercent,
					InodesTotal:       item.InodesTotal,
					InodesUsed:        item.InodesUsed,
					InodesFree:        item.InodesFree,
					InodesUsedPercent: item.InodesUsedPercent,
				},
				Subject: doc.Identity,
			}
			TriggerCheck(doc.Identity, obj, model.DataType_DiskPerformence, doc.Timestamp)
		}
	case map[string]disk.IOCountersStat:
		for disk, item := range val {
			obj := model.PerfData{
				DiskIO: model.DiskIOPerfData{
					Disk:             disk,
					ReadCount:        item.ReadCount,
					MergedReadCount:  item.MergedReadCount,
					WriteCount:       item.WriteCount,
					MergedWriteCount: item.MergedWriteCount,
					ReadBytes:        item.ReadBytes,
					WriteBytes:       item.WriteBytes,
					ReadTime:         item.ReadTime,
					WriteTime:        item.WriteTime,
					IopsInProgress:   item.IopsInProgress,
					IoTime:           item.IoTime,
					WeightedIO:       item.WeightedIO,
					Name:             item.Name,
					SerialNumber:     item.SerialNumber,
					Label:            item.Label,
				},
				Subject: doc.Identity,
			}
			TriggerCheck(doc.Identity, obj, model.DataType_DiskIOPerformence, doc.Timestamp)
		}

	case []net.IOCountersStat:
		for _, item := range val {
			obj := model.PerfData{
				NetDeviceIO: model.NetDeviceIOPerfData{
					Name:        item.Name,
					BytesSent:   item.BytesSent,
					BytesRecv:   item.BytesRecv,
					PacketsSent: item.PacketsSent,
					PacketsRecv: item.PacketsRecv,
					Errin:       item.Errin,
					Errout:      item.Errout,
					Dropin:      item.Dropin,
					Dropout:     item.Dropout,
					Fifoin:      item.Fifoin,
					Fifoout:     item.Fifoout,
				},
				Subject: doc.Identity,
			}
			TriggerCheck(doc.Identity, obj, model.DataType_NetDeviceIOPerformence, doc.Timestamp)
		}
	}
}

func dispatch(arg interface{}) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("get error: ", zap.Error(err.(error)))
		}
	}()
	doc := arg.(*mutual.Document)
	perfData := doc
	values := perfData.Data
	identity := perfData.Identity
	linux := model.GetLinuxIdByIdentity(identity)
	if linux != nil {
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
			saveHostInfo(linux, val, doc.Timestamp)
		case []cpu.InfoStat:
			saveCPUInfo(linux, val, doc.Timestamp)
		case net.InterfaceStatList:
			saveInterfaceInfo(linux, val, doc.Timestamp)
		case mutual.ProcessSnapshot:
			ProcessSnapshotHandler(linux, val, doc.Timestamp)
		}
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

func saveHostInfo(linux *model.Linux, info host.InfoStat, timestamp int64) {
	err := model.UpsertHost(map[string]interface{}{
		"host_identity": linux.LinuxId,
		"info":          info,
		"timestamp":     timestamp,
	})

	if err != nil {
		log.Default().Println("error create linux: ", err)
	}
}

func saveCPUInfo(linux *model.Linux, infoLst []cpu.InfoStat, timestamp int64) {
	err := model.UpdateCPUInfo(map[string]interface{}{"host_identity": linux.LinuxId, "cpu_lst": infoLst, "timestamp": timestamp})

	if err != nil {
		log.Default().Println("error save cpu detail: ", err)
	}
}

func saveInterfaceInfo(linux *model.Linux, infoLst net.InterfaceStatList, timestamp int64) {
	err := model.UpdateInterface(map[string]interface{}{
		"host_identity": linux.LinuxId,
		"interface":     infoLst,
		"timestamp":     timestamp,
	})

	if err != nil {
		log.Default().Println("error save interface detail: ", err)
	}
}

func CachePortPIDMapping(linuxId string, connLst []net.ConnectionStat) map[uint32]int32 {
	key := "port_" + linuxId
	mapping := make(map[uint32]int32)

	var buffer bytes.Buffer
	latest := len(connLst) - 1
	for idx, conn := range connLst {
		buffer.WriteString(strconv.FormatUint(uint64(conn.Laddr.Port), 10))
		buffer.WriteString(":")
		buffer.WriteString(strconv.FormatInt(int64(conn.Pid), 10))
		if latest == idx {
			break
		}
		mapping[conn.Laddr.Port] = conn.Pid
		buffer.WriteString(";")
	}

	model.CacheSet(key, buffer.String(), 12*time.Hour)
	return mapping
}

func isLocalIP(ip string) bool {
	ipObj := gonet.ParseIP(ip)
	return ipObj.IsLoopback()
}

func CacheProcessInfo(linux *model.Linux, procLst []mutual.ProcessInfo, timestamp int64) {
	key := fmt.Sprintf("proc_%d", linux.Id)

	entryLst := make(map[string]any, len(procLst))
	for _, proc := range procLst {
		entryLst[strconv.FormatInt(int64(proc.Pid), 10)] = common.Stringfy(proc)
	}

	entryLst["timestamp"] = strconv.FormatInt(timestamp, 10)

	model.CacheHMSet(key, entryLst)
}

func CacheLocalTCPInfo(linux *model.Linux, connLst []net.ConnectionStat, timestamp int64) {
	for _, conn := range connLst {
		zap.L().Debug("cache connection info",
			zap.String("identity", linux.LinuxId),
			zap.String("ip", conn.Laddr.IP),
			zap.Uint32("port", conn.Laddr.Port),
			zap.Uint32("family", conn.Family),
			zap.Bool("isLocal", isLocalIP(conn.Laddr.IP)),
			zap.Bool("isIPv4", conn.Family == syscall.AF_INET),
		)
		if !isLocalIP(conn.Laddr.IP) && (conn.Family == syscall.AF_INET || conn.Family == syscall.AF_INET6) {
			key := fmt.Sprintf("tcp_%s_%d", conn.Laddr.IP, conn.Laddr.Port)
			value := fmt.Sprintf("%s|%d", linux.LinuxId, conn.Pid)
			model.CacheSet(key, value, 1*time.Hour)
			zap.L().Debug("cache connection info", zap.String("key", key), zap.String("value", value))
		}
	}
}

func GetProcessByIPAndPort(ip string, port uint32) (string, int32) {
	key := fmt.Sprintf("tcp_%s_%d", ip, port)
	value := model.CacheGet(key)

	if value == "" {
		return "", -1
	}

	info := strings.Split(value, "|")
	identity := info[0]
	pid, _ := strconv.ParseInt(info[1], 10, 32)
	return identity, int32(pid)
}

func Transfer2TCPRelation(linux *model.Linux, connLst []net.ConnectionStat, timestamp int64) []*model.TCPRelation {
	result := make([]*model.TCPRelation, 0)
	listening := map[uint32]bool{}

	tcpLst := make([]*net.ConnectionStat, 0)

	for _, conn := range connLst {
		if conn.Status == "LISTEN" {
			listening[conn.Laddr.Port] = true
		} else {
			tcpLst = append(tcpLst, &conn)
		}
	}

	for _, conn := range tcpLst {
		r := new(model.TCPRelation)

		rIp := conn.Raddr.IP
		rPort := conn.Raddr.Port
		zap.L().Debug("got connection", zap.String("ip", rIp), zap.Uint32("port", rPort))

		if listening[conn.Laddr.Port] || (conn.Family != syscall.AF_INET && conn.Family != syscall.AF_INET6) {
			continue
		}

		if isLocalIP(rIp) {
			r.LocalIdentity = linux.LinuxId
			r.LocalPid = conn.Pid
			r.LocalPort = conn.Laddr.Port
			r.LocalIP = conn.Laddr.IP
			r.Timestamp = timestamp

			r.RemotePort = rPort
			r.RemoteIdentity = linux.LinuxId
			r.RemoteIP = conn.Raddr.IP

			for _, c := range connLst {
				if c.Laddr.Port == rPort {
					r.RemotePid = c.Pid
					break
				}
			}
			zap.L().Debug("got local connection", zap.String("ip", rIp), zap.Uint32("port", rPort), zap.Int32("pid", r.RemotePid))
		} else {
			identity, pid := GetProcessByIPAndPort(rIp, rPort)
			if pid == -1 {
				continue
			}

			r.LocalIdentity = linux.LinuxId
			r.LocalPid = conn.Pid
			r.LocalPort = conn.Laddr.Port
			r.LocalIP = conn.Laddr.IP
			r.Timestamp = timestamp

			r.RemotePort = rPort
			r.RemoteIdentity = identity
			r.RemotePid = pid
			r.RemoteIP = conn.Raddr.IP
			zap.L().Debug("got remote connection", zap.String("ip", rIp), zap.Uint32("port", rPort), zap.String("identity", identity), zap.Int32("pid", pid))
		}

		result = append(result, r)
	}
	return result
}

func ProcessListDiff(newLst []mutual.ProcessInfo, oldLst []mutual.ProcessInfo) ([]mutual.ProcessInfo, []mutual.ProcessInfo, []mutual.ProcessInfo) {

	newMapping := make(map[uint32]mutual.ProcessInfo, len(newLst))
	for _, v := range newLst {
		newMapping[v.Hash()] = v
	}

	oldMapping := make(map[uint32]mutual.ProcessInfo, len(oldLst))
	for _, v := range oldLst {
		oldMapping[v.Hash()] = v
	}

	newProcLst := make([]mutual.ProcessInfo, 0)
	expiredLst := make([]mutual.ProcessInfo, 0)
	both := make([]mutual.ProcessInfo, 0)

	for k1, v1 := range newMapping {
		_, exists := oldMapping[k1]
		if !exists {
			newProcLst = append(newProcLst, v1)
		}
	}

	for k1, v1 := range oldMapping {
		_, exists := newMapping[k1]
		if !exists {
			expiredLst = append(expiredLst, v1)
		} else {
			both = append(both, v1)
		}
	}

	return newProcLst, both, expiredLst

}

func ProcessHandler(linux *model.Linux, snapshot mutual.ProcessSnapshot, timestamp int64) map[string]string {

	mapping := make(map[string]string)

	oldProcLst := model.GetProcessLstByLinuxId(linux.LinuxId)
	newProcLst := snapshot.ProcessLst

	newLst, both, expiredLst := ProcessListDiff(newProcLst, oldProcLst)

	if len(expiredLst) > 0 {
		model.RemoveProcInfoFromGraphDB(expiredLst)
	}

	for _, value := range both {
		mapping[fmt.Sprintf("%s|%d", linux.LinuxId, value.Pid)] = value.Id
	}

	if len(newLst) > 0 {
		localMapping, err := model.SaveProcess(linux.LinuxId, newLst, timestamp)
		if err != nil {
			zap.L().Error("error saving process list", zap.Error(err))
			return nil
		}

		for k, v := range localMapping {
			mapping[k] = v
		}

		idLst := make([]string, 0, len(localMapping))
		for _, v := range localMapping {
			idLst = append(idLst, v)
		}

		model.SaveDeploymentRelation(linux.LinuxId, idLst)

	}

	return mapping
}

func ProcessSnapshotHandler(linux *model.Linux, snapshot mutual.ProcessSnapshot, timestamp int64) {

	CacheProcessInfo(linux, snapshot.ProcessLst, timestamp)
	CacheLocalTCPInfo(linux, snapshot.ConnLst, timestamp)

	localMapping := ProcessHandler(linux, snapshot, timestamp)

	zap.L().Debug("got connection list.", zap.Int("length", len(snapshot.ConnLst)))
	relationLst := Transfer2TCPRelation(linux, snapshot.ConnLst, timestamp)
	zap.L().Debug("transfer connection list to relation list.", zap.Int("length", len(relationLst)))

	remoteLst := make([]map[string]any, 0)

	for _, relation := range relationLst {
		relation.From = localMapping[fmt.Sprintf("%s|%d", relation.LocalIdentity, relation.LocalPid)]
		remoteLst = append(remoteLst, map[string]any{
			"pid":      relation.RemotePid,
			"identity": relation.RemoteIdentity,
		})
	}

	remoateMapping := model.QueryProcessIdFromGraphDB(remoteLst)

	for _, relation := range relationLst {
		relation.To = remoateMapping[fmt.Sprintf("%s|%d", relation.RemoteIdentity, relation.RemotePid)]
	}

	model.SaveTCPConnection(RemoveDuplication(relationLst))

}

func RemoveDuplication(rLst []*model.TCPRelation) []*model.TCPRelation {
	unique := make(map[string]bool)
	relationLst := make([]*model.TCPRelation, 0)
	for _, r := range rLst {
		key := fmt.Sprintf("%s|%s", r.From, r.To)
		if _, exists := unique[key]; !exists {
			unique[key] = true
			relationLst = append(relationLst, r)
		}
	}
	return relationLst
}
