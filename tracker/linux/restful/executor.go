package restful

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/syspulse/tracker/linux/client"
	"github.com/syspulse/tracker/linux/common"
	"github.com/syspulse/tracker/linux/logging"
	"github.com/syspulse/tracker/linux/task"
	"github.com/syspulse/tracker/linux/task/kernel"
	"github.com/syspulse/tracker/linux/task/network"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/process"
)

func NewExecutor(client *client.Courier) (*Executor, error) {
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
	client *client.Courier
}

func (ws *Executor) init(client *client.Courier) {
	router := gin.Default()
	router.Use(logging.GinLogger(), logging.GinRecovery(true))
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

func (ws *Executor) RunServer(client *client.Courier) {
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
		job := task.Job{}
		err = json.Unmarshal(body, &job)
		if err != nil {
			log.Default().Println(err)
			return
		}
		jobId := job.Id
		category := job.Category

		switch category {
		case "traffic":
			network.CreateCollectingTask(job)
		case "proc_profiling":
			kernel.CreateProfilingTask(job, func() {
				task.UpdateJobStatus(jobId, task.JOB_STATUS_RUNNING)
			}, func(data []*task.EBPFProfilingData) {
				task.SendResult(jobId, data)
			})
		}
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "ok"})
	})
}
