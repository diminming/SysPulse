package handler

import (
	"net/http"
	"strconv"
	"syspulse/model"
	"syspulse/restful/server/response"

	"github.com/gin-gonic/gin"
)

func QueryCpuPerf(linuxId string, start int64, end int64, cpu string) []map[string]interface{} {

	if cpu == "" {
		return model.DBSelect("select `id`, `user`, `system`, `idle`, `nice`, `iowait`, `irq`, `softirq`, `steal`, `guest`, `guestnice`, `timestamp` from perf_cpu where linux_id = ? and cpu= 'cpu-total' and timestamp between ? and ?;", linuxId, start, end)
	} else {
		return model.DBSelect("select `id`, `user`, `system`, `idle`, `nice`, `iowait`, `irq`, `softirq`, `steal`, `guest`, `guestnice`, `timestamp` from perf_cpu where linux_id = ? and cpu= ? and timestamp between ? and ?;", linuxId, cpu, start, end)
	}

}

func QueryMemoryPerf(linuxId string, start int64, end int64) []map[string]interface{} {
	return model.DBSelect("select `id`, `total`, `available`, `free`, `timestamp` from perf_mem where linux_id = ? and timestamp between ? and ?;", linuxId, start, end)
}

func QueryLoadPerf(linuxId string, start int64, end int64) []map[string]interface{} {
	return model.DBSelect("select `id`, `load1`, `load5`, `load15`, `timestamp` from perf_load where linux_id = ? and timestamp between ? and ?;", linuxId, start, end)
}

func QuerySwapUsed(linuxId string, start int64, end int64) []map[string]interface{} {
	lst := model.DBSelect("select `id`, `used`, `timestamp` from perf_swap where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
	return lst
}

func GetPerfLoad(linuxId string, start int64, end int64) []map[string]interface{} {
	return model.DBSelect("select * from perf_load where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
}

func GetPerfAvailableMemory(linuxId string, start int64, end int64) []map[string]interface{} {
	return model.DBSelect("select `id`, `available`, `timestamp` from perf_mem where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
}

func GetPerfDiskThroughput(linuxId string, start int64, end int64) []map[string]interface{} {
	lst := model.DBSelect("select `id`, `readcount`, `writecount`, `name`, `timestamp` from perf_disk_io where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
	for _, item := range lst {
		item["name"] = string(item["name"].([]uint8))
	}
	return lst
}

func GetPerfIfIO(linuxId string, start int64, end int64) []map[string]interface{} {
	lst := model.DBSelect("select `id`, `name`, `bytessent`, `bytesrecv`, `timestamp` from perf_if_io where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
	for _, item := range lst {
		item["name"] = string(item["name"].([]uint8))
	}
	return lst
}

func QueryPerfSwap(linuxId string, start int64, end int64) []map[string]interface{} {
	return model.DBSelect("select id, total, used, free, sin, sout, pgin, pgout, pgfault, timestamp from perf_swap where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
}

func QueryPerfFS(linuxId string, start int64, end int64) []map[string]interface{} {
	lst := model.DBSelect("select `id`, `linux_id`, `path`, `total`, `free`, `used`, `usedpercent`, `inodestotal`, `inodesused`, `inodesfree`, `inodesusedpercent`, `timestamp` from perf_fs_usage where linux_id = ? and timestamp between ? and ?", linuxId, start, end)
	for _, item := range lst {
		item["path"] = string(item["path"].([]uint8))
	}
	return lst
}

func GetPerfDiskIO(linuxId string, start int64, end int64) []map[string]interface{} {
	sql := "select `id`, `readcount`, `mergedreadcount`, `writecount`, `mergedwritecount`, `readbytes`, `writebytes`, `readtime`, `writetime`, `iopsinprogress`, `iotime`, `weightedio`, `name`, `timestamp` from perf_disk_io where `linux_id` = ? and `timestamp` between ? and ?"
	lst := model.DBSelect(sql, linuxId, start, end)
	for _, item := range lst {
		item["name"] = string(item["name"].([]uint8))
	}
	return lst
}

func GetPerfNetIO(linuxId string, start int64, end int64) []map[string]interface{} {
	sql := "select `name`, `bytessent`, `bytesrecv`, `packetssent`, `packetsrecv`, `errin`, `errout`, `dropin`, `dropout`, `fifoin`, `fifoout`, `timestamp` from perf_if_io where `linux_id`=? and `timestamp` between ? and ?"
	lst := model.DBSelect(sql, linuxId, start, end)
	for _, item := range lst {
		item["name"] = string(item["name"].([]uint8))
	}
	return lst
}

func GetNetPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := GetPerfNetIO(linuxId, start, end)
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetDiskPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := GetPerfDiskIO(linuxId, start, end)
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetCpuUsage(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")
	cpu := ""
	if query.Has("cpu") {
		cpu = query.Get("cpu")
	}
	lst := QueryCpuPerf(linuxId, start, end, cpu)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}
func GetLoad1(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	linuxId := query.Get("linuxId")
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "start is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "end is not a number..."})
		return
	}
	lst := GetPerfLoad(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetMemoryAvailable(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	linuxId := query.Get("linuxId")

	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "start is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "end is not a number..."})
		return
	}
	lst := GetPerfAvailableMemory(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetSwapUsed(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	linuxId := query.Get("linuxId")

	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "start is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "end is not a number..."})
		return
	}
	lst := QuerySwapUsed(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetDiskIOCount(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	linuxId := query.Get("linuxId")

	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "start is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "end is not a number..."})
		return
	}
	lst := GetPerfDiskThroughput(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetInterfaceIOCount(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	linuxId := query.Get("linuxId")

	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "start is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "end is not a number..."})
		return
	}
	lst := GetPerfIfIO(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetCpuPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	cpu := ""
	if query.Has("cpu") {
		cpu = query.Get("cpu")
	}
	lst := QueryCpuPerf(linuxId, start, end, cpu)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetMemoryPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := QueryMemoryPerf(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetLoadPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := QueryLoadPerf(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetSwapPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := QueryPerfSwap(linuxId, start, end)
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}

func GetFSPerf(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	start, err := strconv.ParseInt(query.Get("start"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	end, err := strconv.ParseInt(query.Get("end"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number..."})
		return
	}
	linuxId := query.Get("linuxId")

	lst := QueryPerfFS(linuxId, start, end)
	defer func() {
		lst = lst[:0]
	}()
	ctx.JSON(http.StatusOK, response.JsonResponse{Data: lst, Msg: "success", Status: http.StatusOK})
}
