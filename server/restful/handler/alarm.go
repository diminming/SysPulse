package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
)

func GetAlarmByPage(page int, pageSize int) []model.Alarm {
	sql := "SELECT \n" +
		"    a.id, a.timestamp, linux.hostname, a.ack, a.`trigger`, a.create_timestamp\n" +
		"FROM\n" +
		"    (SELECT \n" +
		"        `id`,\n" +
		"            `timestamp`,\n" +
		"            `linux_id`,\n" +
		"            `trigger`,\n" +
		"            `ack`,\n" +
		"            `perf_data`,\n" +
		"            `create_timestamp`\n" +
		"    FROM\n" +
		"        alarm\n" +
		"    ORDER BY `timestamp` DESC\n" +
		"    LIMIT ? , ?) a\n" +
		"        LEFT JOIN\n" +
		"    linux ON a.linux_id = linux.id;"
	first := page * pageSize
	result := make([]model.Alarm, 0)
	lst := model.DBSelect(sql, first, pageSize)
	for _, item := range lst {
		result = append(result, model.Alarm{
			Id:              item["id"].(int64),
			Timestamp:       item["timestamp"].(int64),
			CreateTimestamp: item["create_timestamp"].(int64),
			Trigger:         string(item["trigger"].([]uint8)),
			Ack:             item["ack"].(int64) == 1,
			Linux: model.Linux{
				Hostname: string(item["hostname"].([]uint8)),
			},
		})
	}
	return result
}

func GetAlarmLstByPage(ctx *gin.Context) {
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
	result := GetAlarmByPage(page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
		"lst":   result,
		"total": model.GetTotalofAlarm(),
	}, Msg: "success"})
}

func GetAlarmById(alarmId int64) *model.Alarm {
	sql := "select a.`id`, a.`timestamp`, a.`linux_id`, a.`trigger`, a.`ack`, a.`perf_data`, a.`create_timestamp`, l.`id` as linuxId, l.`hostname` from (select * from alarm where id = ?) a inner join linux l where l.id = a.linux_id"
	alarmInfo := model.DBSelectRow(sql, alarmId)
	perfData := model.PerfData{}
	err := json.Unmarshal(alarmInfo["perf_data"].([]uint8), &perfData)
	if err != nil {
		log.Default().Printf("error get alarm info at method GetAlarmById: %v\n", err)
	}
	return &model.Alarm{
		Id:              alarmInfo["id"].(int64),
		Timestamp:       alarmInfo["timestamp"].(int64),
		CreateTimestamp: alarmInfo["create_timestamp"].(int64),
		Trigger:         string(alarmInfo["trigger"].([]uint8)),
		Ack:             alarmInfo["ack"].(int64) == 1,
		PerfData:        perfData,
		Linux: model.Linux{
			Id:       alarmInfo["linuxId"].(int64),
			Hostname: string(alarmInfo["hostname"].([]uint8)),
		},
	}
}

func GetAlarm(ctx *gin.Context) {
	alarmId, err := strconv.ParseInt(ctx.Param("alarmId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
	}

	log.Default().Printf("alarm id: %v", alarmId)
	alarm := GetAlarmById(alarmId)

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK", Data: alarm})

}
