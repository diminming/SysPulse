package model

import (
	"database/sql"

	"go.uber.org/zap"
)

type NMON struct {
	Id              int64  `json:"id"`
	Hostname        string `json:"hostname"`
	From            int64  `json:"from"`
	To              int64  `json:"to"`
	Source          string `json:"source"`
	Path            string `json:"path"`
	CreateTimestamp int64  `json:"createTimestamp"`
}

func InsertNMONRecord(nmon *NMON) {
	sql := "INSERT INTO nmon(`hostname`, `from`, `to`, `source`, `path`, `createTimestamp`) VALUES(?, ?, ?, ?, ?, ?)"
	nmon.Id = DBInsert(sql, nmon.Hostname, nmon.From, nmon.To, nmon.Source, nmon.Path, nmon.CreateTimestamp)
}

func GetNMONRecordById(id int64) *NMON {
	nmon := new(NMON)
	sql := "SELECT `id`, `hostname`, `from`, `to`, `source`, `path`, `createTimestamp` FROM nmon WHERE id = ?"
	result := DBSelectRow(sql, id)
	nmon.Id = result["id"].(int64)
	nmon.Path = string(result["path"].([]uint8))
	nmon.Source = string(result["source"].([]uint8))
	nmon.Hostname = string(result["hostname"].([]uint8))
	nmon.From = result["from"].(int64)
	nmon.To = result["to"].(int64)
	nmon.CreateTimestamp = result["createTimestamp"].(int64)
	return nmon
}

func GetNMONRecordByPage(id4Linux int64, page, size int) []*NMON {
	nmonLst := make([]*NMON, 0)
	sql := "SELECT \n" +
		"    n.`id` as id,\n" +
		"    `hostname`,\n" +
		"    `from`,\n" +
		"    `to`,\n" +
		"    `source`,\n" +
		"    `path`,\n" +
		"    `createTimestamp`\n" +
		"FROM\n" +
		"    (SELECT \n" +
		"        id, linux_id\n" +
		"    FROM\n" +
		"        linux\n" +
		"    WHERE\n" +
		"        id = ?) AS l\n" +
		"        LEFT JOIN\n" +
		"    nmon AS n ON l.linux_id = n.hostname\n" +
		"ORDER BY n.createTimestamp DESC , n.id DESC\n" +
		"LIMIT ? , ?;"
	result := DBSelect(sql, id4Linux, page*size, size)
	for _, row := range result {
		id, exists := row["id"]
		if !exists || id == nil {
			continue
		}
		nmon := new(NMON)
		nmon.Id = id.(int64)
		nmon.Path = string(row["path"].([]uint8))
		nmon.Source = string(row["source"].([]uint8))
		nmon.Hostname = string(row["hostname"].([]uint8))
		nmon.From = row["from"].(int64)
		nmon.To = row["to"].(int64)
		nmon.CreateTimestamp = row["createTimestamp"].(int64)
		nmonLst = append(nmonLst, nmon)
	}
	return nmonLst

}

func GetNMONTotal() int64 {
	sqlstr := "select count(id) as total from nmon"
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(sqlstr)
	err := row.Scan(&count)
	if err != nil {
		zap.L().Error("GetNMONTotal", zap.Error(err))
	}
	return count
}
