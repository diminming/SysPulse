package model

import (
	"database/sql"

	"go.uber.org/zap"
)

// 1, "已创建"
// 2, "运行中"
// 3, "已完成"
// 100, "已超时"

const (
	JOB_STATUS_CREATED  = 1
	JOB_STATUS_RUNNING  = 2
	JOB_STATUS_FINISHED = 3
)

type Job struct {
	Id              int64    `json:"id"`
	JobName         string   `json:"job_name"`
	Category        string   `json:"category"`
	Type            string   `json:"type"`
	Status          int      `json:"status"`
	StartupTime     int64    `json:"startup_time"`
	LinuxId         int64    `json:"linux_id"`
	Pid             int32    `json:"pid"`
	Duration        int32    `json:"duration"`
	Immediately     bool     `json:"immediately"`
	IfName          string   `json:"ifName"`
	IpAddr          string   `json:"ipAddr"`
	Port            int32    `json:"port"`
	Direction       []string `json:"direction"`
	Count           int64    `json:"count"`
	Extend          string   `json:"extend"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
}

func GetJobTotal() int64 {
	s := "select count(id) from job"
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		zap.L().Error("error get job total: ", zap.Error(err))
	}
	return count
}

func JobMarkOverdue(timestamp int64) {
	s := "update job set `status` = 100 where (`status` = 1 or `status` = 2) and `create_timestamp` < ?"
	result, err := SqlDB.Exec(s, timestamp)
	if err != nil {
		zap.L().Error("error makr overdue: ", zap.Error(err))
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("error makr overdue: ", zap.Error(err))
		return
	}
	zap.L().Debug("Marked tasks as timed out.", zap.Int64("count", count))
}

func DeleteJob(id int64) int64 {
	s := "delete from job where id = ?"
	result, err := SqlDB.Exec(s, id)
	if err != nil {
		zap.L().Error("error delete job: ", zap.Error(err))
		return -1
	}
	count, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("error delete job: ", zap.Error(err))
		return -1
	}
	zap.L().Debug("delete job record.", zap.Int64("count", count))
	return count
}
