package task

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/syspulse/tracker/linux/client"
	"github.com/syspulse/tracker/linux/common"

	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	EBPFProfilingStackType_PROCESS_KERNEL_SPACE int32 = 0
	EBPFProfilingStackType_PROCESS_USER_SPACE   int32 = 1
	JOB_STATUS_CREATED                          int32 = 1
	JOB_STATUS_RUNNING                          int32 = 2
	JOB_STATUS_FINISHED                         int32 = 3
)

type Job struct {
	Id              int64    `json:"id"`
	JobName         string   `json:"job_name"`
	Category        string   `json:"category"`
	Type            string   `json:"type"`
	Status          int      `json:"status"`
	StartupTime     int64    `json:"startup_time"`
	LinuxId         int64    `json:"linux_id"`
	Pid             int32    `json:"pid"`
	Duration        int32    `json:"duration"`
	Immediately     bool     `json:"immediately"`
	IfName          string   `json:"ifName"`
	IpAddr          string   `json:"ipAddr"`
	Port            int32    `json:"port"`
	Direction       []string `json:"direction"`
	Count           int64    `json:"count"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
}

type RunningSuccessNotify func()

type isEBPFProfilingData_Profiling interface {
	isEBPFProfilingData_Profiling()
}

type EBPFProfilingTaskMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// profiling task id
	TaskId string `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	// profiling process id
	ProcessId string `protobuf:"bytes,2,opt,name=processId,proto3" json:"processId,omitempty"`
	// the start time of this profiling process
	ProfilingStartTime int64 `protobuf:"varint,3,opt,name=profilingStartTime,proto3" json:"profilingStartTime,omitempty"`
	// report time
	CurrentTime int64 `protobuf:"varint,4,opt,name=currentTime,proto3" json:"currentTime,omitempty"`
}

type EBPFProfilingStackMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// stack type
	StackType int32 `protobuf:"varint,1,opt,name=stackType,proto3,enum=skywalking.v3.EBPFProfilingStackType" json:"stackType,omitempty"`
	// stack id from kernel provide
	StackId int32 `protobuf:"varint,2,opt,name=stackId,proto3" json:"stackId,omitempty"`
	// stack symbols
	StackSymbols []string `protobuf:"bytes,3,rep,name=stackSymbols,proto3" json:"stackSymbols,omitempty"`
}

type EBPFProfilingData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// task metadata
	Task *EBPFProfilingTaskMetadata `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
	// profiling data
	//
	// Types that are assignable to Profiling:
	//
	//	*EBPFProfilingData_OnCPU
	//	*EBPFProfilingData_OffCPU
	Profiling isEBPFProfilingData_Profiling `protobuf_oneof:"profiling"`
}

type EBPFOnCPUProfiling struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// stack data in one task(thread)
	Stacks []*EBPFProfilingStackMetadata `protobuf:"bytes,1,rep,name=stacks,proto3" json:"stacks,omitempty"`
	// stack counts
	DumpCount int32 `protobuf:"varint,2,opt,name=dumpCount,proto3" json:"dumpCount,omitempty"`
}

type EBPFProfilingData_OnCPU struct {
	OnCPU *EBPFOnCPUProfiling `protobuf:"bytes,2,opt,name=onCPU,proto3,oneof"`
}

func (*EBPFProfilingData_OnCPU) isEBPFProfilingData_Profiling() {}

func UpdateJobStatus(jobId int64, status int32) {
	srvCfg := common.SysArgs.Server.Restful
	url := fmt.Sprintf("http://%s:%d%s/job/updateStatus", srvCfg.Host, srvCfg.Port, srvCfg.BasePath)
	payload, err := json.Marshal(map[string]interface{}{
		"jobId":  jobId,
		"status": status,
	})
	if err != nil {
		log.Default().Println(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(payload)))
	if err != nil {
		log.Default().Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Println("Error reading response: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Default().Printf("Error updating job status: %d, %s", resp.StatusCode, string(body))
	}
}

func UploadOutcome(fPath string) (string, string) {
	objName := fmt.Sprintf("%s_%s", common.SysArgs.Identity, time.Now().Format("20060102_1504"))
	client.Upload2FileServer("syspulse", objName, fPath, "application/octet-stream")
	return "syspulse", objName
}

func SendResult(jobId int64, data interface{}) {
	srvCfg := common.SysArgs.Server.Restful
	url := fmt.Sprintf("http://%s:%d%s/job/%d/onFinish", srvCfg.Host, srvCfg.Port, srvCfg.BasePath, jobId)
	payload, err := json.Marshal(data)
	if err != nil {
		log.Default().Println(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(payload)))
	if err != nil {
		log.Default().Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Println(err)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Println("Error reading response: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Default().Printf("Error updating job status: %d, %s", resp.StatusCode, string(respBody))
	}
}
