package component

import (
	"log"
	"os"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"gopkg.in/yaml.v2"
)

type Trigger struct {
	LinuxIdentity string   `yaml:"linux"`
	Expres        []string `yaml:"expression"`
}

type TriggerSetting struct {
	TriggerLst []*Trigger `yaml:"triggers"`
}

var TriggerConfig TriggerSetting

var TriggerCache map[string][]*vm.Program

type CpuPerfData struct {
	User      float64 `expr:"user"`
	System    float64 `expr:"system"`
	Idle      float64 `expr:"idle"`
	Nice      float64 `expr:"nice"`
	Iowait    float64 `expr:"iowait"`
	Irq       float64 `expr:"irq"`
	Softirq   float64 `expr:"softirq"`
	Steal     float64 `expr:"steal"`
	Guest     float64 `expr:"guest"`
	GuestNice float64 `expr:"guestNice"`
	CpuUtil   float64 `expr:"util"`
}

type LoadPerfData struct {
	Load1  float64 `expr:"load1"`
	Load5  float64 `expr:"load5"`
	Load15 float64 `expr:"load15"`
}

type MemoryPerfData struct {
	Total          uint64 `expr:"total"`
	Free           uint64 `expr:"free"`
	Active         uint64 `expr:"active"`
	Inactive       uint64 `expr:"inactive"`
	Wired          uint64 `expr:"wired"`
	Laundry        uint64 `expr:"laundry"`
	Buffers        uint64 `expr:"buffers"`
	Cached         uint64 `expr:"cached"`
	WriteBack      uint64 `expr:"writeBack"`
	Dirty          uint64 `expr:"dirty"`
	WriteBackTmp   uint64 `expr:"writeBackTmp"`
	Shared         uint64 `expr:"shared"`
	Slab           uint64 `expr:"slab"`
	Sreclaimable   uint64 `expr:"sreclaimable"`
	Sunreclaim     uint64 `expr:"sunreclaim"`
	PageTables     uint64 `expr:"pageTables"`
	SwapCached     uint64 `expr:"swapCached"`
	CommitLimit    uint64 `expr:"commitLimit"`
	CommittedAS    uint64 `expr:"committedAS"`
	HighTotal      uint64 `expr:"highTotal"`
	HighFree       uint64 `expr:"highFree"`
	LowTotal       uint64 `expr:"lowTotal"`
	LowFree        uint64 `expr:"lowFree"`
	SwapTotal      uint64 `expr:"swapTotal"`
	SwapFree       uint64 `expr:"swapFree"`
	Mapped         uint64 `expr:"mapped"`
	VmallocTotal   uint64 `expr:"vmallocTotal"`
	VmallocUsed    uint64 `expr:"vmallocUsed"`
	VmallocChunk   uint64 `expr:"vmallocChunk"`
	HugePagesTotal uint64 `expr:"hugePagesTotal"`
	HugePagesFree  uint64 `expr:"hugePagesFree"`
	HugePagesRsvd  uint64 `expr:"hugePagesRsvd"`
	HugePagesSurp  uint64 `expr:"hugePagesSurp"`
	HugePageSize   uint64 `expr:"hugePageSize"`
	AnonHugePages  uint64 `expr:"anonHugePages"`
}

type SwapPerfData struct {
	Total       uint64  `expr:"total"`
	Used        uint64  `expr:"used"`
	Free        uint64  `expr:"free"`
	UsedPercent float64 `expr:"usedPercent"`
	Sin         uint64  `expr:"sin"`
	Sout        uint64  `expr:"sout"`
	PgIn        uint64  `expr:"pgIn"`
	PgOut       uint64  `expr:"pgOut"`
	PgFault     uint64  `expr:"pgFault"`
	PgMajFault  uint64  `expr:"pgMajFault"`
}

type DiskPerfData struct {
	Path              string  `expr:"path"`
	Fstype            string  `expr:"fstype"`
	Total             uint64  `expr:"total"`
	Free              uint64  `expr:"free"`
	Used              uint64  `expr:"used"`
	UsedPercent       float64 `expr:"usedPercent"`
	InodesTotal       uint64  `expr:"inodesTotal"`
	InodesUsed        uint64  `expr:"inodesUsed"`
	InodesFree        uint64  `expr:"inodesFree"`
	InodesUsedPercent float64 `expr:"inodesUsedPercent"`
}
type DiskIOPerfData struct {
	Disk             string `expr:"disk"`
	ReadCount        uint64 `expr:"readCount"`
	MergedReadCount  uint64 `expr:"mergedReadCount"`
	WriteCount       uint64 `expr:"writeCount"`
	MergedWriteCount uint64 `expr:"mergedWriteCount"`
	ReadBytes        uint64 `expr:"readBytes"`
	WriteBytes       uint64 `expr:"writeBytes"`
	ReadTime         uint64 `expr:"readTime"`
	WriteTime        uint64 `expr:"writeTime"`
	IopsInProgress   uint64 `expr:"iopsInProgress"`
	IoTime           uint64 `expr:"ioTime"`
	WeightedIO       uint64 `expr:"weightedIO"`
	Name             string `expr:"name"`
	SerialNumber     string `expr:"serialNumber"`
	Label            string `expr:"label"`
}

type NetDeviceIOPerfData struct {
	Name        string `expr:"name"`
	BytesSent   uint64 `expr:"bytesSent"`
	BytesRecv   uint64 `expr:"bytesRecv"`
	PacketsSent uint64 `expr:"packetsSent"`
	PacketsRecv uint64 `expr:"packetsRecv"`
	Errin       uint64 `expr:"errin"`
	Errout      uint64 `expr:"errout"`
	Dropin      uint64 `expr:"dropin"`
	Dropout     uint64 `expr:"dropout"`
	Fifoin      uint64 `expr:"fifoin"`
	Fifoout     uint64 `expr:"fifoout"`
}

type PerfData struct {
	CPU         CpuPerfData         `expr:"cpu"`
	Load        LoadPerfData        `expr:"load"`
	Memory      MemoryPerfData      `expr:"memory"`
	Swap        SwapPerfData        `expr:"swap"`
	Disk        DiskPerfData        `expr:"disk"`
	DiskIO      DiskIOPerfData      `expr:"diskio"`
	NetDeviceIO NetDeviceIOPerfData `expr:"netio"`
}

func init() {

	yamlFile, err := os.ReadFile(common.SysArgs.TriggerCfg)
	if err != nil {
		log.Default().Fatalf("can't open config file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &TriggerConfig)
	if err != nil {
		log.Fatalf("can't read config file: %v", err)
	}
	TriggerCache = make(map[string][]*vm.Program)
	buildTrigger()
}

func buildTrigger() {

	for _, item := range TriggerConfig.TriggerLst {
		identity := item.LinuxIdentity
		expressions := item.Expres
		expObjLst := make([]*vm.Program, 0, len(expressions))
		for _, expression := range expressions {
			expObj, err := expr.Compile(expression, expr.Env(PerfData{}), expr.AsBool())
			if err != nil {
				log.Default().Fatalf("can't create expression by %s.\n%v\n", expression, err)
			}
			expObjLst = append(expObjLst, expObj)
		}
		TriggerCache[identity] = expObjLst
	}
}

func createAlarmRecord(timestamp int64, identity string, trigger string, parameters PerfData) {
	linuxId := model.CacheGet(identity)
	perfDataStr := common.ToString(parameters)

	sql := "insert into alarm(`timestamp`,`linux_id`,`trigger`,`ack`,`perf_data`,`create_timestamp`) value(?,?,?,?,?,?)"

	model.DBInsert(sql, timestamp, linuxId, trigger, false, []byte(perfDataStr), time.Now().UnixMilli())
}

func TriggerCheck(identity string, parameters PerfData, timestamp int64) {
	programs := TriggerCache[identity]
	for _, program := range programs {

		result, err := expr.Run(program, parameters)

		if err != nil {
			log.Default().Panicf("error calc at linux: %s, exp: %s, data: %s. \n%v", identity, program.Source().String(), common.ToString(parameters), err)
		}

		if result.(bool) {
			log.Default().Printf("<Alarm> timestamp: %d, identity: %s, result: %b, trigger: %s", timestamp, identity, result, program.Source().String())
			createAlarmRecord(timestamp, identity, program.Source().String(), parameters)
		}

	}
}
