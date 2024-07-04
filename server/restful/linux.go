package restful

import (
	"encoding/json"
	"fmt"
	"insight/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetLinuxByPage(page int, pageSize int) []model.Linux {
	first := page * pageSize
	sql := `SELECT 
    a.id, 
    a.hostname, 
    a.linux_id, 
    a.biz_id,
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
    LIMIT ? , ?) a
        LEFT JOIN
    biz b ON a.biz_id = b.id`
	lst := model.DBSelect(sql, first, pageSize)
	result := []model.Linux{}
	for _, o := range lst {

		var biz model.Business

		if o["biz_id"].(int64) == 0 {
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

func CreateLinux(linux *model.Linux) {
	sql := "insert into linux(`hostname`, `linux_id`, `biz_id`, `agent_conn`, create_timestamp, update_timestamp) value(?, ?, ?, ?, ?)"
	id := model.DBInsert(sql, linux.Hostname, linux.LinuxId, linux.Biz.Id, linux.AgentConn, linux.CreateTimestamp, linux.UpdateTimestamp)
	linux.Id = id
}

func UpdateLinux(linux *model.Linux, id int64) {
	sql := "update linux set `id` = ?, `hostname`=?, `linux_id`=?, `biz_id`=?, `agent_conn`=?, `update_timestamp`=? where `id`=?"
	model.DBUpdate(sql, linux.Id, linux.Hostname, linux.LinuxId, linux.Biz.Id, linux.AgentConn, linux.UpdateTimestamp, id)
}

func DeleteLinux(linuxId int) {
	sql := "delete from linux where id = ?"
	model.DBDelete(sql, linuxId)
}

func GetLinuxById(id int64) *model.Linux {
	sql := "select * from linux where id = ?"
	target := model.DBSelectRow(sql, id)
	linux := new(model.Linux)
	if target["agent_conn"] != nil {
		linux.AgentConn = string(target["agent_conn"].([]uint8))
	}
	// linux.AgentConn = string(target["agent_conn"].([]uint8))
	linux.CreateTimestamp = target["create_timestamp"].(int64)
	linux.UpdateTimestamp = target["update_timestamp"].(int64)
	linux.Id = target["id"].(int64)
	linux.Hostname = string(target["hostname"].([]uint8))
	linux.LinuxId = string(target["linux_id"].([]uint8))
	return linux
}

func (ws *WebServer) MappingHandler4Linux() {
	ws.Get("/linux/:id", func(ctx *gin.Context) {
		idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
		}
		linux := GetLinuxById(int64(idOfLinux))
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: linux, Msg: "success"})
	})
	ws.Get("/linux/page", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()
		page, err := strconv.Atoi(values.Get("page"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
			return
		}
		pageSize, err := strconv.Atoi(values.Get("pageSize"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
			return
		}
		lst := GetLinuxByPage(page, pageSize)
		total := GetLinuxTotal()
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
			"lst":   lst,
			"total": total,
		}, Msg: "success"})
	})

	ws.Post("/linux", func(ctx *gin.Context) {
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
		CreateLinux(&linux)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: &linux, Msg: "success"})
	})

	ws.Put("/linux/:id", func(ctx *gin.Context) {
		idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
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
		UpdateLinux(linux, idOfLinux)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: &linux, Msg: "success"})
	})

	ws.Delete("/linux", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()
		linuxId, err := strconv.Atoi(values.Get("linux_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
			return
		}
		DeleteLinux(linuxId)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "success"})
	})

	ws.Get("/linux/:id/procLst", func(ctx *gin.Context) {
		idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
			return
		}
		resp := GetProcLst(int64(idOfLinux))
		procLst := resp["data"]
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: procLst, Msg: "success"})
	})

	ws.Get("/linux/:id/proc/:pid/analyze", func(ctx *gin.Context) {
		idOfLinux, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "linux id is not a number."})
			return
		}
		pid, err := strconv.ParseInt(ctx.Param("pid"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "process id is not a number."})
			return
		}
		jobLst := GetAnalyzationJobLst(int64(idOfLinux), pid)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: jobLst, Msg: "success"})
	})
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

func GetProcLst(id int64) map[string]interface{} {

	linux := GetLinuxById(id)
	agentConn := linux.AgentConn

	resp, err := http.Get(fmt.Sprintf("http://%s/api/proc/lst", agentConn))
	if err != nil {
		log.Default().Print("Error sending request: ", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Print("Error reading response: ", err)
		return nil
	}
	if 200 != resp.StatusCode {
		panic(fmt.Sprintf("Error getting process list: %d, %s", resp.StatusCode, string(body)))
	}
	procLst := make(map[string]interface{}, 0)

	json.Unmarshal(body, &procLst)

	return procLst
}
