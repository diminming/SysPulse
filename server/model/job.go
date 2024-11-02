package model

import (
	"database/sql"

	"go.uber.org/zap"
)

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
	s := "update job set `status` = 100 where `status` != 100 and `create_timestamp` < ?"
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
	zap.L().Info("Marked tasks as timed out.", zap.Int64("count", count))
}
