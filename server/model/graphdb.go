package model

import (
	"context"
	"log"
	"time"

	"github.com/syspulse/common"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// var graphClient driver.Client

var (
	graphClient driver.Client
	graphDB     driver.Database
	Endpoints   = common.SysArgs.Storage.GraphDB.Endpoints
	Username    = common.SysArgs.Storage.GraphDB.Username
	Password    = common.SysArgs.Storage.GraphDB.Password
	DBName      = common.SysArgs.Storage.GraphDB.DBName
)

func initConn() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: Endpoints,
	})
	if err != nil {
		panic(err)
	}
	newClient, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(Username, Password),
	})
	if err != nil {
		panic(err)
	}
	graphClient = newClient
}

func initDB() {
	ctx := context.Background()
	// 检查数据库是否存在
	dbExists, err := graphClient.DatabaseExists(ctx, DBName)
	if err != nil {
		log.Fatalf("error check DB: %v", err)
	}
	db, err := graphClient.Database(ctx, DBName)
	if err != nil {
		log.Fatalf("error get DB: %v", err)
	}
	graphDB = db

	// 如果数据库不存在，创建它
	if !dbExists {
		// 创建数据库
		db, err := graphClient.CreateDatabase(ctx, DBName, nil)
		if err != nil {
			log.Fatalf("error create graph db: %v", err)
		}
		graphDB = db
	}
}

var LST_COLLECTION_SETTING = []map[string]string{
	{
		"name": "business",
		"type": "doc",
	},
	{
		"name": "host",
		"type": "doc",
	},
	{
		"name": "process",
		"type": "doc",
	},
	{
		"name": "conn_tcp",
		"type": "edge",
	},
	{
		"name": "deployment",
		"type": "edge",
	},
	{
		"name": "res_consumption",
		"type": "edge",
	},
}

func initCollection() {
	ctx := context.Background()
	for _, setting := range LST_COLLECTION_SETTING {
		collName := setting["name"]
		exists, err := graphDB.CollectionExists(ctx, collName)
		if err != nil {
			log.Fatalf("error check collection: %v", err)
		}
		if !exists {
			switch setting["type"] {
			case "doc":
				_, err := graphDB.CreateCollection(ctx, collName, &driver.CreateCollectionOptions{
					Type: driver.CollectionTypeDocument, // 设置为 Document 类型
				})
				if err != nil {
					log.Fatalf("error create collection '%s': %v", collName, err)
				}
			case "edge":
				_, err := graphDB.CreateCollection(ctx, collName, &driver.CreateCollectionOptions{
					Type: driver.CollectionTypeEdge, // 设置为 Document 类型
				})
				if err != nil {
					log.Fatalf("error create collection '%s': %v", collName, err)
				}
			}
		}
	}
}

func initGraph() {
	ctx := context.Background()
	exists, err := graphDB.GraphExists(ctx, "graph_demployment")
	if err != nil {
		log.Fatalf("error check graph 'graph_demployment', %v", err)
	}
	if !exists {
		_, err := graphDB.CreateGraph(ctx, "graph_demployment", &driver.CreateGraphOptions{
			EdgeDefinitions: []driver.EdgeDefinition{
				{
					Collection: "conn_tcp",
					From:       []string{"process"},
					To:         []string{"process"},
				},
				{
					Collection: "deployment",
					From:       []string{"process"},
					To:         []string{"host"},
				},
				{
					Collection: "res_consumption",
					From:       []string{"business"},
					To:         []string{"host"},
				},
			},
		})
		if err != nil {
			log.Fatalf("error create graph 'graph_demployment', %v", err)
		}
	}
}

func init() {
	initConn()
	initDB()
	initCollection()
	initGraph()
}

func UpdateCPUInfo(doc interface{}) error {
	ctx := context.Background()
	aql := `LET doc = ` + common.ToString(doc) + `
FOR h IN host
FILTER h.host_identity == doc.host_identity
UPDATE h with doc IN host`
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{})
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
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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

	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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
	_, err := graphDB.Query(ctx, aql, map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		log.Fatalf("Failed to execute transaction: %v", err)
	}
}

// func GetLinuxIdByIP(ip string) ([]int64, error) {
// 	aql := `
// for h in host
//   for i in h.interface
//     for addr in i.addrs
//       filter addr.addr like @ip
//         return h.host_identity
// `
// 	ctx := context.Background()
// 	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
// 		"ip": ip + "/%",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cur.Close()

// 	keyLst := make([]int64, 0)
// 	for {
// 		var linuxId int64
// 		_, err := cur.ReadDocument(context.Background(), &linuxId)
// 		if driver.IsNoMoreDocuments(err) {
// 			break
// 		} else if err != nil {
// 			return nil, err
// 		}
// 		keyLst = append(keyLst, linuxId)
// 	}
// 	return keyLst, nil
// }

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
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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
	// DeleteTimeoutHostRecord(timestamp)
	DeleteTimeoutProcessRecord(timestamp)
	DeleteTimeoutTCPRecord(timestamp)
	DeleteTimeoutDeploymentRecord(timestamp)
}

// func DeleteTimeoutHostRecord(timestamp int64) {
// 	aql := `
// FOR t IN host
// 	FILTER t.timestamp < @timestamp
// 	REMOVE { _key: t._key } IN host
// `
// 	ctx := context.Background()
// 	GraphDB.Query(ctx, aql, map[string]interface{}{
// 		"timestamp": timestamp,
// 	})
// }

func DeleteTimeoutProcessRecord(timestamp int64) {
	aql := `
FOR t IN process
	FILTER t.timestamp < @timestamp
	REMOVE { _key: t._key } IN process
`
	ctx := context.Background()
	graphDB.Query(ctx, aql, map[string]interface{}{
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
	graphDB.Query(ctx, aql, map[string]interface{}{
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
	graphDB.Query(ctx, aql, map[string]interface{}{
		"timestamp": timestamp,
	})
}

func QueryLinuxDesc(id int64) map[string]any {
	aql := `
for h in host
  filter h.host_identity == @id 
  return {
    "base": h.info,
    "ifLst": h.interface,
    "cpu": h.cpu_lst
  }
`
	ctx := context.Background()
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		log.Default().Panicln("error load linux description: ", err)
	}
	defer cur.Close()

	result := make(map[string]interface{})
	_, err = cur.ReadDocument(context.Background(), &result)
	if err != nil {
		log.Default().Panicln("error load linux description: ", err)
	}

	return result
}

func QueryLinuxTopo(linuxId int64) ([]map[string]interface{}, error) {
	aql := `
for h in host
  filter h.host_identity == @linuxId
  for v, e, p in 1..2 any h graph graph_demployment
  return p
`
	ctx := context.Background()
	cur, err := graphDB.Query(ctx, aql, map[string]interface{}{
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
		cur, err := graphDB.Query(ctx, aql, map[string]any{
			"offset": times * batchSize,
			"size":   batchSize,
		})

		if err != nil {
			log.Default().Println("error in BatchGetNIC: ", err)
			return
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

func SaveBiz(biz *Business) {
	ctx := context.Background()
	coll, err := graphDB.Collection(ctx, "business")
	if err != nil {
		log.Default().Println("can't get business collection: ", err)
	}
	_, err = coll.CreateDocument(ctx, biz)
	if err != nil {
		log.Default().Println("error create business document: ", err)
	}
}

func SaveConsumptionRelation(linux *Linux) {
	ctx := context.Background()
	aql := `
for biz in business
  filter biz.id == @bizId
  limit 1
  for linux in host
    filter linux.host_identity == @linuxId
    limit 1
  let _from = concat("business/", biz._key)
  let _to = concat("host/", linux._key)
  
  insert
  {
    "_from": _from,
    "_to": _to,
    "timestamp": @timestamp
  }
  into res_consumption
`
	_, err := graphDB.Query(ctx, aql, map[string]any{
		"bizId":     linux.Biz.Id,
		"linuxId":   linux.Id,
		"timestamp": time.Now().UnixMilli(),
	})
	if err != nil {
		log.Default().Printf("Failed to execute transaction: %v", err)
	}

}

func UpdateConsumptionRelation(linux *Linux) {
	ctx := context.Background()
	aql := `FOR biz IN business
  FILTER biz.id == @bizId
  LIMIT 1
  FOR linux IN host
    FILTER linux.host_identity == @linuxId
    LIMIT 1
	UPDATE linux._key with {"name": @name} IN host
    UPSERT {"_to": linux._id} 
      INSERT {
        "timestamp": @timestamp,
        "_from": biz._id, 
        "_to": linux._id
      } 
      UPDATE {
        "_from": biz._id, 
        "timestamp": @timestamp,
      }  
      IN res_consumption`

	meta, err := graphDB.Query(ctx, aql, map[string]any{
		"bizId":     linux.Biz.Id,
		"linuxId":   linux.Id,
		"name":      linux.LinuxId,
		"timestamp": time.Now().UnixMilli(),
	})

	if err != nil {
		log.Default().Printf("Failed to execute transaction: %v", err)
	}

	log.Default().Println("UpdateConsumptionRelation: ", meta)
}
