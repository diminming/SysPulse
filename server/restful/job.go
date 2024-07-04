package restful

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syspulse/common"
	"syspulse/model"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateJob(job *model.Job) (*model.Job, error) {

	tx, err := model.SqlDB.Begin()
	if err != nil {
		log.Default().Println(err)
	}
	defer tx.Rollback()
	sql := "insert into job(`job_name`, `category`, `type`, `status`, `startup_time`, `linux_id`, `pid`, `duration`, `immediately`, `create_timestamp`, `update_timestamp`) value(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(sql, job.JobName, job.Category, job.Type, job.Status, job.StartupTime, job.LinuxId, job.Pid, job.Duration, job.Immediately, job.CreateTimestamp, job.UpdateTimestamp)
	if err != nil {
		log.Default().Println(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Default().Println(err)
	}
	job.Id = insertId

	linuxId := job.LinuxId
	linux := GetLinuxById(linuxId)

	reqBody, err := json.Marshal(job)
	if err != nil {
		log.Default().Print("Error converting job info: ", err)
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/job", linux.AgentConn), "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		log.Default().Print("Error sending request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Print("Error reading response: ", err)
		return nil, err
	}

	if 200 != resp.StatusCode {
		panic(fmt.Sprintf("Error creating job: %d, %s", resp.StatusCode, string(body)))
	}

	err = tx.Commit()
	if err != nil {
		log.Default().Println(err)
	}

	return job, nil

}

func Write2File(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 使用缓冲区写入数据
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	println("Binary data has been written to", filename)
}

func ReadFromFile(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		log.Default().Println(err)
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data := make([]byte, 0, 64*1024)
	// 创建一个字节切片用于存储数据
	buffer := make([]byte, 1024) // 选择合适的缓冲区大小

	// 从文件中读取数据
	for {
		// 从 Reader 中读取数据到缓冲区
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		// 处理读取的数据
		data = append(data, buffer[:n]...)
	}

	return string(data), nil

}

func (ws *WebServer) MappingHandler4Job() {
	ws.Post("/job", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Default().Println(err)
			return
		}
		job := new(model.Job)
		err = json.Unmarshal(body, job)
		if err != nil {
			log.Default().Println(err)
			return
		}
		job.Status = model.JOB_STATUS_CREATED
		job.CreateTimestamp = time.Now().Unix()
		job.UpdateTimestamp = time.Now().Unix()
		CreateJob(job)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: job, Msg: "success"})
	})
	ws.PatchWithGrp("/job/updateStatus", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Default().Println(err)
			return
		}
		data := new(map[string]interface{})
		err = json.Unmarshal(body, data)
		if err != nil {
			log.Default().Println(err)
			return
		}
		jobId := int64((*data)["jobId"].(float64))
		status := int32((*data)["status"].(float64))
		model.DBUpdate("update job set `status` = ? where `id` = ?", status, jobId)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "ok"})
	}, ws.CallbackGrp)

	ws.PatchWithGrp("/job/:jobId/onFinish", func(ctx *gin.Context) {

		jobId, err := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
		}

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Default().Println(err)
			return
		}
		data := new([]map[string]interface{})
		err = json.Unmarshal(body, data)
		if err != nil {
			log.Default().Println(err)
			return
		}
		filename := fmt.Sprintf("%s/insight_%d.json", common.SysArgs.Storage.File.Path, jobId)

		Write2File(filename, body)

		model.DBUpdate("update job set `status` = ? where `id` = ?", model.JOB_STATUS_FINISHED, jobId)

		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "ok"})
	}, ws.CallbackGrp)

	ws.Get("/job/:jobId", func(ctx *gin.Context) {
		jobId, err := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
			return
		}
		filename := fmt.Sprintf("%s/insight_%d.json", common.SysArgs.Storage.File.Path, jobId)
		ctn, err := ReadFromFile(filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
			return
		}

		var data []map[string]interface{}

		// 解析JSON字符串到map
		err = json.Unmarshal([]byte(ctn), &data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, JsonResponse{Status: http.StatusInternalServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "ok", Data: data})
	})
}
