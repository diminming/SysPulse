package server

import (
	"net/http"
	"syspulse/restful/handler"
)

func (ws *WebServer) SetupRoutes() {
	ws.Register4Api(http.MethodGet, "/linux/:id", handler.GetLinuxInfoById)
	ws.Register4Api(http.MethodGet, "/linux/page", handler.GetLinuxLstByPage)
	ws.Register4Api(http.MethodPost, "/linux", handler.NewLinuxRecord)
	ws.Register4Api(http.MethodPut, "/linux/:id", handler.UpdateLinuxRecord)
	ws.Register4Api(http.MethodDelete, "/linux", handler.DeleteLinuxRecord)
	ws.Register4Api(http.MethodGet, "/linux/:id/procLst", handler.GetProcessLst)
	ws.Register4Api(http.MethodGet, "/linux/:id/proc/:pid/analyze", handler.GetProcAnalJobLst)

	ws.Register4Api(http.MethodPost, "/job", handler.CreateJob)
	ws.Register4Api(http.MethodGet, "/job/:jobId", handler.GetJobResult)
	ws.Register4Callback(http.MethodPatch, "/job/updateStatus", handler.UpdateJobStatus)
	ws.Register4Callback(http.MethodPatch, "/job/:jobId/onFinish", handler.OnJobFinish)

	ws.Register4Api(http.MethodPost, "/biz", handler.NewBizRecord)
	ws.Register4Api(http.MethodPut, "/biz", handler.UpdateBizRecord)
	ws.Register4Api(http.MethodGet, "/biz/page", handler.GetBizLstByPage)
	ws.Register4Api(http.MethodDelete, "/biz", handler.DeleteBizRecord)

	ws.Register4Api(http.MethodGet, "/cache/page", handler.GetCacheRecordLstByPage)
	ws.Register4Api(http.MethodPost, "/cache", handler.NewCacheRecord)
	ws.Register4Api(http.MethodPut, "/cache", handler.UpdateCacheRecord)
	ws.Register4Api(http.MethodDelete, "/cache", handler.DeleteCacheRecord)

	ws.Register4Api(http.MethodPost, "/login", handler.UserLogin)

	ws.Register4Api(http.MethodGet, "/db/page", handler.GetDBRecordLstByPage)
	ws.Register4Api(http.MethodPost, "/db", handler.NewDBRecord)
	ws.Register4Api(http.MethodPut, "/db", handler.ModifyDBRecord)
	ws.Register4Api(http.MethodDelete, "/db", handler.RemoveDBRecord)

	ws.Register4Api(http.MethodGet, "/perf/net", handler.GetNetPerf)
	ws.Register4Api(http.MethodGet, "/perf/disk", handler.GetDiskPerf)
	ws.Register4Api(http.MethodGet, "/perf/cpu/usage", handler.GetCpuUsage)
	ws.Register4Api(http.MethodGet, "/perf/load/load1", handler.GetLoad1)
	ws.Register4Api(http.MethodGet, "/perf/mem/available", handler.GetMemoryAvailable)
	ws.Register4Api(http.MethodGet, "/perf/swap/used", handler.GetSwapUsed)
	ws.Register4Api(http.MethodGet, "/perf/disk/iocount", handler.GetDiskIOCount)
	ws.Register4Api(http.MethodGet, "/perf/if/iocount", handler.GetInterfaceIOCount)
	ws.Register4Api(http.MethodGet, "/perf/cpu", handler.GetCpuPerf)
	ws.Register4Api(http.MethodGet, "/perf/mem", handler.GetMemoryPerf)
	ws.Register4Api(http.MethodGet, "/perf/load", handler.GetLoadPerf)
	ws.Register4Api(http.MethodGet, "/perf/swap", handler.GetSwapPerf)
	ws.Register4Api(http.MethodGet, "/perf/fs", handler.GetFSPerf)
}
