package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
)

func GetAlarmByPage(page int, pageSize int) []model.Alarm {
	sql := "SELECT \n" +
		"    a.id, a.timestamp, linux.hostname, a.ack, a.`msg`, a.create_timestamp\n" +
		"FROM\n" +
		"    (SELECT \n" +
		"        `id`,\n" +
		"            `timestamp`,\n" +
		"            `linux_id`,\n" +
		"            `msg`,\n" +
		"            `ack`,\n" +
		"            `create_timestamp`\n" +
		"    FROM\n" +
		"        alarm\n" +
		"    ORDER BY `timestamp` DESC\n" +
		"    LIMIT ? , ?) a\n" +
		"        INNER JOIN\n" +
		"    linux ON a.linux_id = linux.id;"
	first := page * pageSize
	result := make([]model.Alarm, 0)
	lst := model.DBSelect(sql, first, pageSize)
	for _, item := range lst {

		result = append(result, model.Alarm{
			Id:              item["id"].(int64),
			Timestamp:       item["timestamp"].(int64),
			CreateTimestamp: item["create_timestamp"].(int64),
			Msg:             string(item["msg"].([]uint8)),
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

func GetAlarmById(alarmId int64) map[string]any {
	sql := "select a.`id`, a.`timestamp`, a.`linux_id`, a.`trigger`, a.`trigger_id`, a.`msg`, a.`ack`, a.`create_timestamp`, l.`id` as linuxId, l.`hostname`, l.`linux_id` from (select * from alarm where id = ?) a inner join linux l where l.id = a.linux_id"
	alarmInfo := model.DBSelectRow(sql, alarmId)

	return map[string]any{
		"id":              alarmInfo["id"].(int64),
		"timestamp":       alarmInfo["timestamp"].(int64),
		"createTimestamp": alarmInfo["create_timestamp"].(int64),
		"trigger":         string(alarmInfo["trigger"].([]uint8)),
		"triggerId":       string(alarmInfo["trigger_id"].([]uint8)),
		"msg":             string(alarmInfo["msg"].([]uint8)),
		"ack":             alarmInfo["ack"].(int64) == 1,
		"linux": map[string]any{
			"id":       alarmInfo["linuxId"].(int64),
			"hostname": string(alarmInfo["hostname"].([]uint8)),
			"linuxId":  string(alarmInfo["linux_id"].([]uint8)),
		},
	}
}

func GetAlarmCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK", Data: model.GetTotalofActiveAlarm()})
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

func GetData4AlarmHeatMap(from, to int64) []map[string]any {
	sqlstr := "select a.time_tag, a.total, biz.`biz_name` from (select biz_id, time_tag, count(id) as total from alarm where `timestamp` between ? and ? group by biz_id, time_tag order by total desc limit 10) a left join biz on biz.id = a.biz_id"
	result := make([]map[string]any, 0, 10)
	lst := model.DBSelect(sqlstr, from, to)

	for _, item := range lst {

		bizName := item["biz_name"]
		if bizName == nil {
			bizName = ""
		} else {
			bizName = string(item["biz_name"].([]uint8))
		}

		total, _ := item["total"].(int64)

		result = append(result, map[string]any{
			"timetag": string(item["time_tag"].([]uint8)),
			"bizName": bizName,
			"total":   total,
		})
	}

	return result
}

func Stat4HeatMap(ctx *gin.Context) {
	from, _ := strconv.ParseInt(ctx.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(ctx.Query("to"), 10, 64)
	result := GetData4AlarmHeatMap(from, to)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK", Data: result})
}

func GetData4AlarmTrend(from, to int64) []map[string]any {
	sqlstr := "select time_tag, count(id) as total from alarm where `timestamp` between ? and ? group by time_tag"
	result := make([]map[string]any, 0, 10)
	lst := model.DBSelect(sqlstr, from, to)

	for _, item := range lst {
		total, _ := item["total"].(int64)
		result = append(result, map[string]any{
			"timetag": string(item["time_tag"].([]uint8)),
			"total":   total,
		})
	}

	return result
}

func Stat4Trend(ctx *gin.Context) {
	from, _ := strconv.ParseInt(ctx.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(ctx.Query("to"), 10, 64)
	result := GetData4AlarmTrend(from, to)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK", Data: result})
}

func UpdateAlarmStatusInDB(id int64, status bool) {
	sqlstr := "update alarm set ack = ? where id = ?"
	_, err := model.DBUpdate(sqlstr, status, id)
	if err != nil {
		log.Default().Panicln("error disable alarm, id: ", id, err)
	}
}

func UpdateAlarmStatusInCache(alarmInfo map[string]any) {
	linux := alarmInfo["linux"].(map[string]any)
	identity := linux["linuxId"].(string)
	triggerId := alarmInfo["triggerId"].(string)
	key := "alarm_" + identity

	model.CacheHSet(key, triggerId, "false")
}

func UpdateAlarmStatus(id int64, status bool) {
	alarmInfo := GetAlarmById(id)

	UpdateAlarmStatusInCache(alarmInfo)
	UpdateAlarmStatusInDB(id, status)
}

func DisableAlarm(ctx *gin.Context) {
	alarmId, err := strconv.ParseInt(ctx.Param("alarmId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
	}
	UpdateAlarmStatus(alarmId, true)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK"})
}
