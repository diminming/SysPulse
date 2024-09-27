package model

import (
	"database/sql"
	"log"
)

func GetJobTotal() int64 {
	s := "select count(id) from job"
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}
