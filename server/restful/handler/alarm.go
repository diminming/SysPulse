package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
)

func genSqlWhere(status string, from, util int64, targetLst []int64) ([]string, []any) {
	conditionLst := make([]string, 0)
	args := make([]any, 0)

	if status != "" {
		conditionLst = append(conditionLst, "ack = ?")
		if status == "active" {
			args = append(args, 0)
		} else {
			args = append(args, 1)
		}
	}

	if from > 0 && util > 0 {
		conditionLst = append(conditionLst, "timestamp between ? and ?")
		args = append(args, from)
		args = append(args, util)
	}

	if len(targetLst) > 0 {
		condition_in := new(strings.Builder)
		condition_in.WriteString("linux_id in (")
		for idx, id := range targetLst {
			if idx > 0 {
				condition_in.WriteString(", ")
			}
			condition_in.WriteString("?")
			args = append(args, id)
		}
		condition_in.WriteString(")")
		conditionLst = append(conditionLst, condition_in.String())
	}

	return conditionLst, args
}

func GetAlarmByPage(status string, from, util int64, targetLst []int64, page int, pageSize int) []model.Alarm {
	first := page * pageSize
	result := make([]model.Alarm, 0)

	sqlArgs := make([]any, 0)
	sqlstr := new(strings.Builder)
	sqlstr.WriteString("SELECT a.id, a.timestamp, linux.hostname, a.ack, a.`msg`, a.create_timestamp ")
	sqlstr.WriteString("FROM ")
	sqlstr.WriteString("(SELECT  ")
	sqlstr.WriteString("`id`, ")
	sqlstr.WriteString("`timestamp`, ")
	sqlstr.WriteString("`linux_id`, ")
	sqlstr.WriteString("`msg`, ")
	sqlstr.WriteString("`ack`, ")
	sqlstr.WriteString("`create_timestamp` ")
	sqlstr.WriteString("FROM ")
	sqlstr.WriteString("alarm ")

	conditionLst, args := genSqlWhere(status, from, util, targetLst)
	if len(conditionLst) > 0 {
		sqlstr.WriteString("WHERE ")
		sqlstr.WriteString(strings.Join(conditionLst, " and "))
		sqlstr.WriteString(" ")
		sqlArgs = append(sqlArgs, args...)
	}

	sqlstr.WriteString("ORDER BY `timestamp` DESC ")
	sqlstr.WriteString("LIMIT ? , ?) a ")
	sqlArgs = append(sqlArgs, first)
	sqlArgs = append(sqlArgs, pageSize)
	sqlstr.WriteString("INNER JOIN ")
	sqlstr.WriteString("linux ON a.linux_id = linux.id;")

	lst := model.DBSelect(sqlstr.String(), sqlArgs...)
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

func GetTotalofAlarm(status string, from, util int64, targetLst []int64) int64 {
	sqlstr := new(strings.Builder)
	sqlstr.WriteString("select count(id) as total from alarm ")
	sqlArgs := make([]any, 0)
	conditionLst, args := genSqlWhere(status, from, util, targetLst)
	if len(conditionLst) > 0 {
		sqlstr.WriteString("WHERE ")
		sqlstr.WriteString(strings.Join(conditionLst, " and "))
		sqlArgs = append(sqlArgs, args...)
	}
	result := model.DBSelectRow(sqlstr.String(), sqlArgs...)
	switch v := result["total"].(type) {
	case int64:
		return v
	case []uint8:
		count, _ := strconv.ParseInt(string(v), 10, 64)
		return count
	}
	return -1
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
	status := values.Get("status")
	fromStr := values.Get("from")
	from := int64(0)
	if fromStr != "" {
		from, err = strconv.ParseInt(fromStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "from is not a number."})
			return
		}
	}

	utilStr := values.Get("util")
	util := int64(0)
	if utilStr != "" {
		util, err = strconv.ParseInt(values.Get("util"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "util is not a number."})
			return
		}
	}

	target := values.Get("target")
	targetLst := make([]int64, 0)
	if target != "" {
		for _, idStr := range strings.Split(target, ",") {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "target lst contain string can't convert to int64."})
				return
			}
			targetLst = append(targetLst, id)
		}
	}
	result := GetAlarmByPage(status, from, util, targetLst, page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
		"lst":   result,
		"total": GetTotalofAlarm(status, from, util, targetLst),
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
