package model

import (
	"database/sql"
	"log"
)

type PerformenceDataType int

const (
	DataType_CpuPerformence         = 1
	DataType_LoadPerformence        = 2
	DataType_MemoryPerformence      = 3
	DataType_SwapPerformence        = 4
	DataType_DiskPerformence        = 5
	DataType_DiskIOPerformence      = 6
	DataType_NetDeviceIOPerformence = 7
)

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

type Alarm struct {
	Id              int64    `json:"id"`
	Timestamp       int64    `json:"timestamp"`
	CreateTimestamp int64    `json:"createTimestamp"`
	TriggerId       string   `json:"triggerId"`
	Trigger         string   `json:"trigger"`
	Ack             bool     `json:"ack"`
	Msg             string   `json:"msg"`
	Linux           Linux    `json:"linux"`
	PerfData        PerfData `json:"perfData"`
}

func GetTotalofAlarm() uint32 {
	s := "select count(id) from alarm"
	var row *sql.Row
	var count uint32
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Default().Printf("%v", err)
	}
	return count
}
