package model

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"

	driver "github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type CIDRRecord struct {
	CIDR    string `json:"cidr"`
	LinuxId string `json:"linuxId"`
}

type Linux struct {
	Id              int64    `json:"id"`
	Hostname        string   `json:"hostname"`
	LinuxId         string   `json:"linux_id"`
	Biz             Business `json:"biz"`
	AgentConn       string   `json:"agent_conn"`
	ExtId           string   `json:"ext_id"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
}

func LoadLinuxByIdentity(identity string) *Linux {
	sqlstr := "select `id`, `hostname`, `linux_id`, `biz_id`, `agent_conn`, `ext_id`, `create_timestamp`, `update_timestamp` from `linux` where linux_id = ?"
	row := DBSelectRow(sqlstr, identity)
	linux := new(Linux)
	linux.Id = row["id"].(int64)
	linux.Hostname = string(row["hostname"].([]uint8))
	linux.LinuxId = string(row["linux_id"].([]uint8))
	linux.AgentConn = string(row["agent_conn"].([]uint8))
	linux.Biz.Id = row["biz_id"].(int64)
	linux.CreateTimestamp = row["create_timestamp"].(int64)
	linux.UpdateTimestamp = row["update_timestamp"].(int64)

	if extId, ok := row["ext_id"]; ok && extId != nil {
		linux.ExtId = string(extId.([]uint8))
	}

	return linux
}

func GetLinuxTotal(keyword string) int64 {
	sqlbuf := new(strings.Builder)
	sqlArgs := make([]any, 0)
	sqlbuf.WriteString("select count(id) from linux ")
	if keyword != "" && !(strings.TrimSpace(keyword) == "") {
		sqlbuf.WriteString("where hostname like ? or linux_id like ?\n")
		likeArg := "%" + keyword + "%"
		sqlArgs = append(sqlArgs, likeArg)
		sqlArgs = append(sqlArgs, likeArg)
	}
	var row *sql.Row
	var count int64
	row = SqlDB.QueryRow(sqlbuf.String(), sqlArgs...)
	err := row.Scan(&count)
	if err != nil {
		log.Default().Println(err)
	}
	return count
}

func GetLinuxIdByIdentity(identity string) *Linux {

	linuxId := CacheGet(identity)
	if linuxId == "0" || linuxId == "" {
		sqlstr := "select id from linux where linux_id = ?"
		// var row *sql.Row
		linux := new(Linux)
		row := SqlDB.QueryRow(sqlstr, identity)
		err := row.Scan(&linux.Id)
		if err != nil {
			zap.L().Error("can't get linux record from sqldb with linuxId: ", zap.String("linuxId", identity))
			return nil
		}
		SetIdentityAndIdMappingInCache(linux)
		linux.LinuxId = identity
		return linux
	} else {
		linux := new(Linux)
		num, err := strconv.ParseInt(linuxId, 10, 64)
		if err != nil {
			log.Default().Println(err)
			return nil
		}
		linux.Id = num
		linux.LinuxId = identity
		return linux
	}
}

func GetInterfaceLst(id string) ([]map[string]interface{}, error) {
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
