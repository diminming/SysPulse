package model

import (
	"database/sql"
	"log"
)

type Business struct {
	Id              int64  `json:"id"`
	BizName         string `json:"bizName"`
	BizId           string `json:"bizId"`
	BizDesc         string `json:"bizDesc"`
	CreateTimestamp int64  `json:"createTimestamp"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
}

func GetBizTotal() int64 {
	s := "select count(id) from biz"
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}
