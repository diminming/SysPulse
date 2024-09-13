package handler

import (
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
