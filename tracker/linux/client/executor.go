package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"syspulse/tracker/linux/common"
	"syspulse/tracker/linux/task"
	"syspulse/tracker/linux/task/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/process"
)

func NewExecutor(client *Courier) (*Executor, error) {
	executor := new(Executor)
	executor.init(client)
	executor.mapping()
	return executor, nil
}

type JsonResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

type Executor struct {
	router *gin.Engine
	appSrv *gin.RouterGroup
	client *Courier
}

func (ws *Executor) init(client *Courier) {
	router := gin.Default()
	appSrv := router.Group(common.SysArgs.Restful.BasePath)

	// appSrv.Use(func(c *gin.Context) {
	// 	path := c.FullPath()

	// 	token := c.Request.Header.Get("token")

	// 	if (len(token) < 1 || len(model.CacheGet(token)) < 1) && path != fmt.Sprintf("%s/login", common.SysArgs.Server.BasePath) {
	// 		c.JSON(http.StatusForbidden, JsonResponse{Status: http.StatusForbidden, Msg: "non-logging"})
	// 		c.Abort()
	// 		return
	// 	} else {
	// 		model.CacheExpire(token, time.Minute*15)
	// 	}

	// 	c.Next()
	// })

	ws.router = router
	ws.appSrv = appSrv
	ws.client = client
}

func (ws *Executor) Post(url string, handler gin.HandlerFunc) {
	ws.appSrv.POST(url, handler)
}

func (ws *Executor) Put(url string, handler gin.HandlerFunc) {
	ws.appSrv.PUT(url, handler)
}

func (ws *Executor) Get(url string, handler gin.HandlerFunc) {
	ws.appSrv.GET(url, handler)
}

func (ws *Executor) Delete(url string, handler gin.HandlerFunc) {
	ws.appSrv.DELETE(url, handler)
}

func (ws *Executor) RunServer(client *Courier) {
	ws.router.Run(common.SysArgs.Restful.Addr)
}

func (ws *Executor) Close() {
	log.Default().Println("Executor is stoped.")
}

func (ws *Executor) mapping() {
	ws.Get("/proc/lst", func(ctx *gin.Context) {
		procLst := []map[string]interface{}{}
		processes, err := process.Processes()
		if err != nil {
			log.Default().Println("Error getting process list: ", err)
		}
		for _, item := range processes {
			proc := make(map[string]interface{})
			proc["pid"] = item.Pid
			proc["name"], _ = item.Name()
			proc["ppid"], _ = item.Ppid()
			proc["create_time"], _ = item.CreateTime()
			proc["exec"], _ = item.Exe()
			procLst = append(procLst, proc)
		}
		ctx.JSON(http.StatusOK, JsonResponse{Data: map[string]interface{}{
			"timestamp": time.Now().UnixMilli(),
			"procLst":   procLst,
		}, Status: http.StatusOK, Msg: "ok"})
	})

	ws.Post("/job", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Default().Println(err)
			return
		}
		jobInfo := new(map[string]interface{})
		err = json.Unmarshal(body, jobInfo)
		if err != nil {
			log.Default().Println(err)
			return
		}
		jobId := int64((*jobInfo)["id"].(float64))
		kernel.CreateProfilingTask(*jobInfo, func() {
			UpdateJobStatus(jobId, task.JOB_STATUS_RUNNING)
		}, func(data []*task.EBPFProfilingData) {
			FinishJob(jobId, data)
		})
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "ok"})
	})
}

func FinishJob(jobId int64, data []*task.EBPFProfilingData) {
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
