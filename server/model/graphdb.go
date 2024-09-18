package model

import (
	"context"
	"log"
	"time"

	"github.com/syspulse/common"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var Client driver.Client
var GraphDB driver.Database

var (
	Endpoints = common.SysArgs.Storage.GraphDB.Endpoints
	Username  = common.SysArgs.Storage.GraphDB.Username
	Password  = common.SysArgs.Storage.GraphDB.Password
	DbName    = common.SysArgs.Storage.GraphDB.DbName
)

func init() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: Endpoints,
	})
	if err != nil {
		panic(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(Username, Password),
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := client.Database(ctx, DbName)
	if err != nil {
		panic(err)
	}
	Client = client
	GraphDB = db
}

func UpdateCPUInfo(doc interface{}) error {
	ctx := context.Background()
	aql := `LET doc = ` + common.ToString(doc) + `
FOR h IN host
FILTER h.host_identity == doc.host_identity
UPDATE h with doc IN host`
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{})
	if err != nil {
		return err
	}
	defer cur.Close()
	return nil
}

func UpsertHost(doc interface{}) error {
	ctx := context.Background()
	aql := `LET doc = @doc
UPSERT {"host_identity": doc.host_identity} 
    INSERT doc 
	UPDATE doc 
	IN host
RETURN NEW._key
`
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return err
	}
	defer cur.Close()
	return nil
}

func UpsertProcess(doc interface{}) ([]string, error) {
	ctx := context.Background()
	aql := `
LET doc = @doc
UPSERT {"host_identity": doc.host_identity, "pid": doc.pid} 
	INSERT doc 
	UPDATE doc 
	IN process
RETURN NEW._key
`

	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	keyLst := make([]string, 0)
	for {
		var key string
		_, err := cur.ReadDocument(context.Background(), &key)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		keyLst = append(keyLst, key)
	}
	return keyLst, nil
}

func UpdateInterface(doc interface{}) error {
	ctx := context.Background()
	aql := `LET doc = @doc
FOR h IN host
FILTER h.host_identity == doc.host_identity
UPDATE h with doc IN host`
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return err
	}
	defer cur.Close()
	return nil
}

func UpsertDeploymentRelation(doc map[string]interface{}) {
	aql := `
LET doc = @doc
LET procLst = doc.procLst
FOR h IN host
	FILTER h.host_identity == doc.host_identity
	LIMIT 1
	for proc in procLst
		LET from = CONCAT("process/", proc)
		LET to = CONCAT("host/", h._key)
		UPSERT {"_from": from, "_to": to} 
			INSERT {
				"timestamp": doc.timestamp,
				"_from": from, 
				"_to": to
			} 
			UPDATE {
				"timestamp": doc.timestamp,
			}  
			IN deployment
`
	ctx := context.Background()
	_, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		log.Fatalf("Failed to execute transaction: %v", err)
	}
}

func GetLinuxIdByIP(ip string) ([]int64, error) {
	aql := `
for h in host
  for i in h.interface
    for addr in i.addrs
      filter addr.addr like @ip
        return h.host_identity
`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"ip": ip + "/%",
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	keyLst := make([]int64, 0)
	for {
		var linuxId int64
		_, err := cur.ReadDocument(context.Background(), &linuxId)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		keyLst = append(keyLst, linuxId)
	}
	return keyLst, nil
}

func UpsertConnRelation(localLinuxId int64, localPid int32, remoteLinuxId int64, remotePid int32, timestamp int64) error {
	aql := `
LET _from = (
  FOR p IN process
    FILTER p.host_identity == @localId && p.pid == @localPid
    LIMIT 1
    RETURN p
)

LET _to = (
  FOR p IN process
    FILTER p.host_identity == @remoteId && p.pid == @remotePid
    LIMIT 1
    RETURN p
)

UPSERT {
  "_from": _from[0]._id,
  "_to": _to[0]._id
}
INSERT {
  "_from": _from[0]._id,
  "_to": _to[0]._id,
  "timestamp": @timestamp
}
UPDATE {
  "timestamp": @timestamp
}
IN conn_tcp
`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"localId":   localLinuxId,
		"localPid":  localPid,
		"remoteId":  remoteLinuxId,
		"remotePid": remotePid,
		"timestamp": timestamp,
	})
	if err != nil {
		return err
	}
	defer cur.Close()
	return nil
}

func DeleteTimeoutTopo(timestamp int64) {
	DeleteTimeoutHostRecord(timestamp)
	DeleteTimeoutProcessRecord(timestamp)
	DeleteTimeoutTCPRecord(timestamp)
	DeleteTimeoutDeploymentRecord(timestamp)
}

func DeleteTimeoutHostRecord(timestamp int64) {
	aql := `
FOR t IN host
	FILTER t.timestamp < @timestamp
	REMOVE { _key: t._key } IN host
`
	ctx := context.Background()
	GraphDB.Query(ctx, aql, map[string]interface{}{
		"timestamp": timestamp,
	})
}

func DeleteTimeoutProcessRecord(timestamp int64) {
	aql := `
FOR t IN process
	FILTER t.timestamp < @timestamp
	REMOVE { _key: t._key } IN process
`
	ctx := context.Background()
	GraphDB.Query(ctx, aql, map[string]interface{}{
		"timestamp": timestamp,
	})
}

func DeleteTimeoutTCPRecord(timestamp int64) {
	aql := `
FOR t IN conn_tcp
	FILTER t.timestamp < @timestamp
	REMOVE { _key: t._key } IN conn_tcp
`
	ctx := context.Background()
	GraphDB.Query(ctx, aql, map[string]interface{}{
		"timestamp": timestamp,
	})
}
func DeleteTimeoutDeploymentRecord(timestamp int64) {
	aql := `
FOR t IN deployment
	FILTER t.timestamp < @timestamp
	REMOVE { _key: t._key } IN deployment
`
	ctx := context.Background()
	GraphDB.Query(ctx, aql, map[string]interface{}{
		"timestamp": timestamp,
	})
}

func QueryLinuxTopo(linuxId int64) ([]map[string]interface{}, error) {
	aql := `
for h in host
  filter h.host_identity == @linuxId
  for v, e, p in 1..2 any h graph graph_demployment
  return p
`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"linuxId": linuxId,
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

func BatchGetNIC(callback func([]map[string]any) bool) {
	times := 0
	batchSize := 100
	startTimestamp := time.Now().UnixMilli()
	for {
		aql := `
for h in host
  for i in h.interface
    for addr in i.addrs
		limit @offset, @size
		return {
			"addr": addr.addr,
			"host_id": h.host_identity
		}
`
		result := false
		ctx := context.Background()
		cur, err := GraphDB.Query(ctx, aql, map[string]any{
			"offset": times * batchSize,
			"size":   batchSize,
		})

		if err != nil {
			log.Default().Println("error in BatchGetNIC: ", err)
		}
		lst := make([]map[string]any, 0)
		for {
			info := make(map[string]interface{})
			_, err := cur.ReadDocument(context.Background(), &info)
			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				log.Default().Println("error in BatchGetNIC: ", err)
			}
			lst = append(lst, info)
		}

		result = callback(lst)
		defer cur.Close()
		times += 1
		if result || (time.Now().UnixMilli()-startTimestamp >= 30*1000) {
			break
		}

		time.Sleep(time.Second)
	}

}
