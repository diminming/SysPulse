package model

import (
	"context"
	"log"
	"syspulse/common"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var Client driver.Client
var GraphDB driver.Database

func init() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		panic(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "123456"),
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := client.Database(ctx, "insight")
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

func UpsertConnRelation(localLinuxId int64, localPid int32, remoteLinuxId int64, remotePid int32, timestamp int64) {
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
		log.Default().Printf("Failed to execute transaction: %v", err)
	}
	defer cur.Close()
}
