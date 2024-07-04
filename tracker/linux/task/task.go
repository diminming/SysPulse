package task

import (
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	EBPFProfilingStackType_PROCESS_KERNEL_SPACE int32 = 0
	EBPFProfilingStackType_PROCESS_USER_SPACE   int32 = 1
	JOB_STATUS_CREATED                          int32 = 1
	JOB_STATUS_RUNNING                          int32 = 2
	JOB_STATUS_FINISHED                         int32 = 3
)

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
