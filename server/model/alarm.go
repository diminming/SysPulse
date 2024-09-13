package model

import (
	"database/sql"
	"log"
)

type Alarm struct {
	Id              int64
	Timestamp       int64
	CreateTimestamp int64
	Trigger         string
	Ack             bool
	Linux           Linux
}

func GetTotalofAlarm() uint32 {
	s := "select count(id) from biz"
	var row *sql.Row
	var count uint32
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Default().Printf("%v", err)
	}
	return count
}
