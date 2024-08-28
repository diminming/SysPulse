package model

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	driver "github.com/arangodb/go-driver"
)

type Linux struct {
	Id              int64    `json:"id"`
	Hostname        string   `json:"hostname"`
	LinuxId         string   `json:"linux_id"`
	Biz             Business `json:"biz"`
	AgentConn       string   `json:"agent_conn"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
}

func GetLinuxTotal() int64 {
	s := "select count(id) from linux"
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(s)
	err := row.Scan(&count)
	if err != nil {
		log.Default().Println(err)
	}
	return count
}

func GetLinuxById(id string) *Linux {

	linuxId := CacheGet(id)
	if linuxId == "0" || linuxId == "" {
		sql := "select id from linux where linux_id = ?"
		// var row *sql.Row
		linux := new(Linux)
		row := SqlDB.QueryRow(sql, id)
		err := row.Scan(&linux.Id)
		if err != nil {
			log.Default().Println(err)
		}
		CacheSet(id, strconv.FormatInt(linux.Id, 10), 0)
		return linux
	} else {
		linux := new(Linux)
		num, err := strconv.ParseInt(linuxId, 10, 64)
		if err != nil {
			log.Default().Println(err)
			return nil
		}
		linux.Id = num
		return linux
	}
}

func GetInterfaceLst(id int64) ([]map[string]interface{}, error) {
	aql := `for h in host
filter h.host_identity == @host_identity
return {
    "if_lst": h.interface
}`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"host_identity": id,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close()
	result := make([]map[string]interface{}, 0, 10)
	for {
		info := make(map[string]interface{})
		_, err := cur.ReadDocument(context.Background(), &info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}
