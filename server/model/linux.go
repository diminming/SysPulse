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

func LoadLinuxByIdentity(identity string) *Linux {
	sqlstr := "select `id`, `hostname`, `linux_id`, `biz_id`, `agent_conn`, `create_timestamp`, `update_timestamp` from `linux` where linux_id = ?"
	row := DBSelectRow(sqlstr, identity)
	linux := new(Linux)
	linux.Id = row["id"].(int64)
	linux.Hostname = string(row["hostname"].([]uint8))
	linux.LinuxId = string(row["linux_id"].([]uint8))
	linux.AgentConn = string(row["agent_conn"].([]uint8))
	linux.Biz.Id = row["biz_id"].(int64)
	linux.CreateTimestamp = row["create_timestamp"].(int64)
	linux.UpdateTimestamp = row["update_timestamp"].(int64)
	return linux
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

func GetLinuxIdByIdentity(id string) *Linux {

	linuxId := CacheGet(id)
	if linuxId == "0" || linuxId == "" {
		sqlstr := "select id from linux where linux_id = ?"
		// var row *sql.Row
		linux := new(Linux)
		row := SqlDB.QueryRow(sqlstr, id)
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
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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
