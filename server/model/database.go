package model

import (
	"database/sql"
	"log"
)

type Database struct {
	Id              int64    `json:"id"`    // record id
	Name            string   `json:"name"`  // database name
	DBID            string   `json:"db_id"` // identifier of DB
	Type            string   `json:"type"`  // database type, MySQL, PGSql, Oracle, DB2...
	Linux           Linux    `json:"linux"`
	Biz             Business `json:"biz"`
	UpdateTimestamp int64    `json:"update_timestamp"`
	CreateTimestamp int64    `json:"create_timestamp"`
}

func GetDBTotal() int {
	s := "select count(1) from db_record"
	var row *sql.Row
	var count int
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Print(err)
	}
	return count
}
