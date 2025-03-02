package server

import (
	"net/http"

	"github.com/syspulse/restful/handler"
)

func (ws *WebServer) SetupRoutes() {
	ws.Register4Api(http.MethodPost, "/biz", handler.NewBizRecord)
	ws.Register4Api(http.MethodPut, "/biz", handler.UpdateBizRecord)
	ws.Register4Api(http.MethodGet, "/biz/page", handler.GetBizLstByPage)
	ws.Register4Api(http.MethodDelete, "/biz", handler.DeleteBizRecord)
	ws.Register4Api(http.MethodGet, "/biz/count", handler.GetBizCount)
	ws.Register4Api(http.MethodGet, "/biz/:bizId", handler.GetBizById)
	ws.Register4Api(http.MethodPost, "/biz/:bizId/topo", handler.QueryBizTopo)
	ws.Register4Api(http.MethodGet, "/biz/:bizId/count_inst", handler.CountInst)

	ws.Register4Api(http.MethodGet, "/linux/:id", handler.GetLinuxInfoById)
	ws.Register4Api(http.MethodGet, "/linux/page", handler.GetLinuxLstByPage)
	ws.Register4Api(http.MethodPost, "/linux", handler.CreateLinuxRecord)
	ws.Register4Api(http.MethodPut, "/linux/:id", handler.ModifyLinuxRecord)
	ws.Register4Api(http.MethodDelete, "/linux", handler.DeleteLinuxRecord)
	ws.Register4Api(http.MethodGet, "/linux/:id/procLst", handler.GetProcessLst)
	ws.Register4Api(http.MethodGet, "/linux/:id/proc/:pid/analyze", handler.GetProcAnalJobLst)
	ws.Register4Api(http.MethodGet, "/linux/:id/topo", handler.GetLinuxTopo)
	ws.Register4Api(http.MethodGet, "/linux/:id/desc", handler.GetLinuxDesc)
	ws.Register4Api(http.MethodGet, "/linux/:id/if/lst", handler.GetInterfaceLst)
	ws.Register4Api(http.MethodGet, "/linux/count", handler.GetLinuxCount)

	ws.Register4Api(http.MethodPost, "/job", handler.CreateJob)
	ws.Register4Api(http.MethodGet, "/job/:jobId/result", handler.GetJobResult)
	ws.Register4Api(http.MethodGet, "/job/traffic/page", handler.GetTrafficJobLst)
	ws.Register4Api(http.MethodGet, "/job/count", handler.GetJobCount)
	ws.Register4Api(http.MethodDelete, "/job/:jobId", handler.DeleteJob)

	ws.Register4Api(http.MethodGet, "/cache/page", handler.GetCacheRecordLstByPage)
	ws.Register4Api(http.MethodPost, "/cache", handler.NewCacheRecord)
	ws.Register4Api(http.MethodPut, "/cache", handler.UpdateCacheRecord)
	ws.Register4Api(http.MethodDelete, "/cache", handler.DeleteCacheRecord)

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

	ws.Register4Callback(http.MethodPatch, "/job/updateStatus", handler.UpdateJobStatus)
	ws.Register4Callback(http.MethodPatch, "/job/:jobId/onFinish", handler.OnJobFinished)

	ws.Register4Api(http.MethodGet, "/alarm/page", handler.GetAlarmLstByPage)
	ws.Register4Api(http.MethodGet, "/alarm/:alarmId", handler.GetAlarm)
	ws.Register4Api(http.MethodPut, "/alarm/:alarmId/disable", handler.DisableAlarm)
	ws.Register4Api(http.MethodGet, "/alarm/count", handler.GetAlarmCount)
	ws.Register4Api(http.MethodGet, "/alarm/stat_heat", handler.Stat4HeatMap)
	ws.Register4Api(http.MethodGet, "/alarm/stat_trend", handler.Stat4Trend)

	ws.Register4Api(http.MethodPost, "/alarm", handler.NewAlarmRecord)

	// 权限相关
	ws.Register4Api(http.MethodPost, "/login", handler.UserLogin)

	ws.Register4Api(http.MethodPost, "/user", handler.CreateUser)
	ws.Register4Api(http.MethodGet, "/user/page", handler.GetUserPage)
	ws.Register4Api(http.MethodDelete, "/user/:userId", handler.DeleteUser)

	ws.Register4Api(http.MethodPost, "/role", handler.CreateRole)
	ws.Register4Api(http.MethodGet, "/role/page", handler.GetRoleByPage)
	ws.Register4Api(http.MethodDelete, "/role/:roleId", handler.DeleteRole)

	ws.Register4Api(http.MethodGet, "/permission/page", handler.GetPermissionLst)
	ws.Register4Api(http.MethodPost, "/permission", handler.CreatePermission)
	ws.Register4Api(http.MethodDelete, "/permission/:prmId", handler.DeletePermission)

	ws.Register4Api(http.MethodPost, "/menu", handler.CreateMenu)
	ws.Register4Api(http.MethodDelete, "/menu/:menuId", handler.DeleteMenu)
	ws.Register4Api(http.MethodGet, "/menu/page", handler.GetMenuByPage)

	ws.Register4Api(http.MethodPost, "/nmon", handler.NewNMONRecord)
	ws.Register4Api(http.MethodGet, "/nmon/:id", handler.GetNMONRecordById)
	ws.Register4Api(http.MethodGet, "/nmon/page", handler.GetNMONRecordByPage)

	ws.Register4Api(http.MethodGet, "/linuxGraph", handler.GetLinuxGraph)
}
