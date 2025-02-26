package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syspulse/component"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"
)

func genSqlWhere(status string, from, util int64, targetLst []int64, bizId int64) ([]string, []any) {
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

	if bizId > 0 {
		conditionLst = append(conditionLst, "biz_id = ?")
		args = append(args, bizId)
	}

	return conditionLst, args
}

func GetAlarmByPage(status string, from, util int64, targetLst []int64, bizId int64, page int, pageSize int) []model.Alarm {
	first := page * pageSize
	result := make([]model.Alarm, 0)

	sqlArgs := make([]any, 0)
	sqlstr := new(strings.Builder)
	sqlstr.WriteString("SELECT a.id, a.timestamp, linux.hostname, a.ack, a.source, a.`msg`, a.create_timestamp, b.id as bizid, b.biz_name ")
	sqlstr.WriteString("FROM ")
	sqlstr.WriteString("(SELECT  ")
	sqlstr.WriteString("`id`, ")
	sqlstr.WriteString("`timestamp`, ")
	sqlstr.WriteString("`linux_id`, ")
	sqlstr.WriteString("`msg`, ")
	sqlstr.WriteString("`ack`, ")
	sqlstr.WriteString("`source`, ")
	sqlstr.WriteString("`create_timestamp` ")
	sqlstr.WriteString("FROM ")
	sqlstr.WriteString("alarm ")

	conditionLst, args := genSqlWhere(status, from, util, targetLst, bizId)
	if len(conditionLst) > 0 {
		sqlstr.WriteString("WHERE ")
		sqlstr.WriteString(strings.Join(conditionLst, " and "))
		sqlstr.WriteString(" ")
		sqlArgs = append(sqlArgs, args...)
	}

	sqlstr.WriteString("ORDER BY `timestamp` DESC, `id` DESC ")
	sqlstr.WriteString("LIMIT ? , ?) a ")
	sqlArgs = append(sqlArgs, first)
	sqlArgs = append(sqlArgs, pageSize)
	sqlstr.WriteString("INNER JOIN ")
	sqlstr.WriteString("linux ON a.linux_id = linux.id left join biz b on linux.biz_id = b.id;")

	lst := model.DBSelect(sqlstr.String(), sqlArgs...)
	for _, item := range lst {

		alarm := model.Alarm{
			Id:              item["id"].(int64),
			Timestamp:       item["timestamp"].(int64),
			CreateTimestamp: item["create_timestamp"].(int64),
			Msg:             string(item["msg"].([]uint8)),
			Ack:             item["ack"].(int64) == 1,
			Source:          string(item["source"].([]uint8)),
			Linux: model.Linux{
				Hostname: string(item["hostname"].([]uint8)),
			},
		}

		bizId, exsits := item["bizid"]
		if exsits && bizId != nil {
			bizIdInt := bizId.(int64)
			if bizIdInt > 0 {
				alarm.Biz = model.Business{
					Id:      bizId.(int64),
					BizName: string(item["biz_name"].([]uint8)),
				}
			}
		}

		result = append(result, alarm)

	}
	return result
}

func GetTotalofAlarm(status string, from, util int64, targetLst []int64, bizId int64) int64 {
	sqlstr := new(strings.Builder)
	sqlstr.WriteString("select count(id) as total from alarm ")
	sqlArgs := make([]any, 0)
	conditionLst, args := genSqlWhere(status, from, util, targetLst, bizId)
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

	bizId := int64(0)

	if values.Has("bizId") {
		bizId, err = strconv.ParseInt(values.Get("bizId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "biz id is not a number."})
			return
		}
	}

	result := GetAlarmByPage(status, from, util, targetLst, bizId, page, pageSize)

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
		"lst":   result,
		"total": GetTotalofAlarm(status, from, util, targetLst, bizId),
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
		zap.L().Panic("error disable alarm, id: ", zap.Int64("id", id), zap.Error(err))
	}
}

func UpdateAlarmStatusInCache(alarmInfo map[string]any) bool {
	linux := alarmInfo["linux"].(map[string]any)
	identity := linux["linuxId"].(string)
	triggerId := alarmInfo["triggerId"].(string)
	key := "alarm_" + identity

	return model.CacheHDel(key, triggerId)
}

func UpdateAlarmStatus(id int64, status bool) {
	alarmInfo := GetAlarmById(id)

	success := UpdateAlarmStatusInCache(alarmInfo)
	if success {
		UpdateAlarmStatusInDB(id, status)
	}

}

func DisableAlarm(ctx *gin.Context) {
	alarmId, err := strconv.ParseInt(ctx.Param("alarmId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "job id is not a number."})
	}
	UpdateAlarmStatus(alarmId, true)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "OK"})
}

func FormatData(alarm *model.Alarm) error {
	alarm.CreateTimestamp = time.Now().UnixMilli()
	linux := model.LoadLinuxByIdentity(alarm.Linux.LinuxId)
	if linux == nil {
		zap.L().Error("can't get linux by linux identity.", zap.String("identity", alarm.Linux.LinuxId))
		return errors.New("can't get linux by linux identity: ")
	}
	alarm.Linux = *linux
	return nil
}

func NewAlarmRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read from request body", zap.Error(err))
		return
	}
	zap.L().Info("the alarm string from outside.", zap.String("alarm string", string(body)))
	var alarm = new(model.Alarm)
	err = json.Unmarshal(body, alarm)
	if err != nil {
		zap.L().Error("error unmarshal alarm from webhook", zap.Error(err))
		return
	}

	err = FormatData(alarm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: err.Error()})
		return
	}

	component.CreateAlarmRecord(alarm.Timestamp*1000, &alarm.Linux, alarm.TriggerId, alarm.Trigger, alarm.Msg, alarm.Source, alarm.Level)
	zap.L().Info("the alarm obj from outside.", zap.String("alarm info", string(body)))
}
