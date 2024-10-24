package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"

	"github.com/gin-gonic/gin"
)

func GetLinuxCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: model.GetLinuxTotal(), Msg: "success"})
}

func GetInterfaceLst(ctx *gin.Context) {
	linuxId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
	}
	infoLst, err := model.GetInterfaceLst(linuxId)
	if err != nil {
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError, response.JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
	}
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: infoLst, Msg: "success"})
}

func GetLinuxTopo(ctx *gin.Context) {
	linuxId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
	}
	infoLst, err := model.QueryLinuxTopo(linuxId)
	if err != nil {
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError, response.JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
	}
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: infoLst, Msg: "success"})
}

func GetLinuxDesc(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
	}
	desc := model.QueryLinuxDesc(id)
	if err != nil {
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError, response.JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
	}
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: desc, Msg: "success"})
}

func GetLinuxInfoById(ctx *gin.Context) {
	idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
	}
	linux := GetLinuxById(int64(idOfLinux))
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: linux, Msg: "success"})
}

func GetLinuxLstByPage(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}
	lst := GetLinuxByPage(page, pageSize)
	total := GetLinuxTotal()
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
		"lst":   lst,
		"total": total,
	}, Msg: "success"})
}

func Insert2SqlDB(linux *model.Linux) {
	sql := "insert into linux(`hostname`, `linux_id`, `biz_id`, `agent_conn`, create_timestamp, update_timestamp) value(?, ?, ?, ?, ?, ?)"
	id := model.DBInsert(sql, linux.Hostname, linux.LinuxId, linux.Biz.Id, linux.AgentConn, linux.CreateTimestamp, linux.UpdateTimestamp)
	linux.Id = id
}

func InsertLinuxRecord(linux *model.Linux) {

	if LinuxIdExist(linux.Id, linux.LinuxId) {
		log.Default().Panicf("Linux id: \"%s\" exist", linux.LinuxId)
		return
	}

	Insert2SqlDB(linux)
	model.UpsertHost(map[string]any{
		"host_identity": linux.Id,
		"name":          linux.LinuxId,
		"timestamp":     time.Now().UnixMilli(),
	})
}

func CreateLinuxRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Default().Println(err)
		return
	}
	var linux = model.Linux{}
	err = json.Unmarshal(body, &linux)
	if err != nil {
		log.Default().Println(err)
		return
	}
	linux.CreateTimestamp = time.Now().Unix()
	linux.UpdateTimestamp = time.Now().Unix()

	InsertLinuxRecord(&linux)

	if linux.Biz.Id > 0 {
		model.SaveConsumptionRelation(&linux)
	}
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: &linux, Msg: "success"})
}

func UpdateLinuxInSqlDB(linux *model.Linux) {
	sql := "update linux set `hostname`=?, `linux_id`=?, `biz_id`=?, `agent_conn`=?, `update_timestamp`=? where `id`=?"
	model.DBUpdate(sql, linux.Hostname, linux.LinuxId, linux.Biz.Id, linux.AgentConn, linux.UpdateTimestamp, linux.Id)
}

func UpdateLinuxRecord(linux *model.Linux, id int64) {

	if id != linux.Id {
		log.Default().Panicln("The two records that need to be updated have inconsistent ID values.")
		return
	}

	linux0 := GetLinuxById(id)
	if linux0 == nil {
		log.Default().Panicf("Can't get linux record by id: %d", id)
		return
	}

	if LinuxIdExist(linux.Id, linux.LinuxId) {
		log.Default().Panicf("Linux id: \"%s\" exist", linux.LinuxId)
		return
	}

	UpdateLinuxInSqlDB(linux)
	model.UpdateConsumptionRelation(linux)

}

func LinuxIdExist(id int64, s string) bool {
	sqlstr := "select count(id) as total from linux where id != ? and linux_id = ?"
	result := model.DBSelectRow(sqlstr, id, s)
	return (result["total"]).(int64) > 0
}

func ModifyLinuxRecord(ctx *gin.Context) {
	idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var linux = new(model.Linux)
	err = json.Unmarshal(body, linux)
	if err != nil {
		fmt.Println(err)
		return
	}

	linux.UpdateTimestamp = time.Now().Unix()
	UpdateLinuxRecord(linux, idOfLinux)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: &linux, Msg: "success"})
}

func DeleteLinuxRecord(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	linuxId, err := strconv.Atoi(values.Get("linux_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	DeleteLinux(linuxId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func GetProcessLst(ctx *gin.Context) {
	idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
		return
	}
	refresh, _ := strconv.ParseBool(ctx.Query("refresh"))

	if refresh {
		procLst, timestamp := GetProcLst(idOfLinux)
		UpdateProcCache(idOfLinux, procLst, timestamp)
		ctx.JSON(http.StatusOK, response.JsonResponse{
			Status: http.StatusOK,
			Data: map[string]interface{}{
				"procLst":   procLst,
				"timestmap": timestamp,
			},
			Msg: "success"},
		)
	} else {
		ctx.JSON(http.StatusOK, response.JsonResponse{
			Status: http.StatusOK,
			Data:   LoadProcLst(idOfLinux),
			Msg:    "success"},
		)
	}
}

func GetProcAnalJobLst(ctx *gin.Context) {
	idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
		return
	}
	pid, err := strconv.ParseInt(ctx.Param("pid"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "process id is not a number."})
		return
	}
	jobLst := GetAnalyzationJobLst(int64(idOfLinux), pid)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: jobLst, Msg: "success"})
}

func GetLinuxByPage(page int, pageSize int) []model.Linux {
	first := page * pageSize
	sql := `SELECT 
    a.id, 
    a.hostname, 
    a.linux_id, 
    b.id as biz_id,
    b.biz_name, 
    a.create_timestamp, 
    a.update_timestamp
FROM
    (SELECT 
        id,
            hostname,
            linux_id,
            biz_id,
            create_timestamp,
            update_timestamp
    FROM
        linux
    ORDER BY update_timestamp DESC , id DESC
    LIMIT ?, ?) a
        LEFT JOIN
    biz b ON a.biz_id = b.id`
	lst := model.DBSelect(sql, first, pageSize)
	result := []model.Linux{}
	for _, o := range lst {

		var biz model.Business

		if o["biz_id"] == nil || o["biz_id"].(int64) == 0 {
			biz = model.Business{}
		} else {
			biz = model.Business{
				Id:      o["biz_id"].(int64),
				BizName: string(o["biz_name"].([]uint8)),
			}
		}

		item := model.Linux{
			Id:              o["id"].(int64),
			Hostname:        string(o["hostname"].([]uint8)),
			LinuxId:         string(o["linux_id"].([]uint8)),
			Biz:             biz,
			CreateTimestamp: o["create_timestamp"].(int64),
			UpdateTimestamp: o["update_timestamp"].(int64),
		}
		result = append(result, item)
	}
	return result
}

func GetLinuxTotal() int64 {
	return model.GetLinuxTotal()
}

func DeleteLinux(linuxId int) {
	sql := "delete from linux where id = ?"
	model.DBDelete(sql, linuxId)
}

func GetLinuxById(id int64) *model.Linux {
	sql := "select l.id, l.hostname, l.linux_id as linux_identity, l.agent_conn, b.id as bizId, b.biz_name from linux l left join biz b on l.biz_id = b.id where l.id = ?"
	target := model.DBSelectRow(sql, id)
	if target == nil {
		log.Default().Printf("there is no record with id: %d\n", id)
		return nil
	}
	linux := new(model.Linux)
	if target["agent_conn"] != nil {
		linux.AgentConn = string(target["agent_conn"].([]uint8))
	}

	if target["bizId"] != nil {
		linux.Biz.Id = target["bizId"].(int64)
		linux.Biz.BizName = string(target["biz_name"].([]uint8))
	}

	linux.Id = target["id"].(int64)
	linux.Hostname = string(target["hostname"].([]uint8))
	linux.LinuxId = string(target["linux_identity"].([]uint8))
	return linux
}

func GetAnalyzationJobLst(linuxId int64, pid int64) []map[string]interface{} {
	lst := model.DBSelect("select id, job_name, type, status, startup_time, duration, immediately, create_timestamp from job where `category`='proc_profiling' and pid = ? and linux_id = ? order by id desc", pid, linuxId)
	for _, item := range lst {
		item["job_name"] = string(item["job_name"].([]uint8))
		item["type"] = string(item["type"].([]uint8))
	}
	return lst
}

func AnalyzeProcInLinux(linuxId int64, pid int64) map[string]interface{} {
	linux := GetLinuxById(linuxId)
	agentConn := linux.AgentConn

	resp, err := http.Get(fmt.Sprintf("http://%s/api/proc/%d/analyze", agentConn, pid))
	if err != nil {
		log.Default().Print("Error sending request: ", err)
		return nil
	}
	defer resp.Body.Close()
	return nil
}

func LoadProcLst(id int64) map[string]string {
	return model.CacheHGetAll(fmt.Sprintf("proc_%d", id))
}

func GetProcLst(id int64) ([]interface{}, int64) {

	linux := GetLinuxById(id)
	agentConn := linux.AgentConn

	resp, err := http.Get(fmt.Sprintf("http://%s/api/proc/lst", agentConn))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Error getting process list: %d, %s", resp.StatusCode, string(body)))
	}
	result := map[string]interface{}{}

	json.Unmarshal(body, &result)

	data := result["data"].(map[string]interface{})

	return data["procLst"].([]interface{}), int64(data["timestamp"].(float64))
}

func UpdateProcCache(id int64, procLst []interface{}, timestamp int64) {
	key := fmt.Sprintf("proc_%d", id)

	entryLst := map[string]interface{}{}
	for _, procInfo := range procLst {
		proc := procInfo.(map[string]interface{})
		entryLst[strconv.FormatInt(int64(proc["pid"].(float64)), 10)] = common.ToString(procInfo)
	}
	entryLst["timestamp"] = strconv.FormatInt(timestamp, 10)

	model.CacheHMSet(key, entryLst)
}
