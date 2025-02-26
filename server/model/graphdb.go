package model

import (
	"context"
	"fmt"
	"log"

	"github.com/syspulse/common"
	"github.com/syspulse/mutual"
	"go.uber.org/zap"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// var graphClient driver.Client

type TCPRelation struct {
	From string `json:"from"`
	To   string `json:"to"`

	LocalIdentity string `json:"local_identity"`
	LocalPid      int32  `json:"local_pid"`
	LocalIP       string `json:"local_ip"`
	LocalPort     uint32 `json:"local_port"`

	RemoteIP       string `json:"remote_ip"`
	RemotePort     uint32 `json:"remote_port"`
	RemoteIdentity string `json:"remote_identity"`
	RemotePid      int32  `json:"remote_pid"`

	Timestamp int64 `json:"timestamp"`
}

var (
	graphClient driver.Client
	GraphDB     driver.Database
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

	// 如果数据库不存在，创建它
	if !dbExists {
		// 创建数据库
		db, err := graphClient.CreateDatabase(ctx, DBName, nil)
		if err != nil {
			log.Fatalf("error create graph db: %v", err)
		}
		GraphDB = db
	} else {
		db, err := graphClient.Database(ctx, DBName)
		if err != nil {
			log.Fatalf("error get DB: %v", err)
		}
		GraphDB = db
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
		exists, err := GraphDB.CollectionExists(ctx, collName)
		if err != nil {
			log.Fatalf("error check collection: %v", err)
		}
		if !exists {
			switch setting["type"] {
			case "doc":
				_, err := GraphDB.CreateCollection(ctx, collName, &driver.CreateCollectionOptions{
					Type: driver.CollectionTypeDocument, // 设置为 Document 类型
				})
				if err != nil {
					log.Fatalf("error create collection '%s': %v", collName, err)
				}
			case "edge":
				_, err := GraphDB.CreateCollection(ctx, collName, &driver.CreateCollectionOptions{
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
	exists, err := GraphDB.GraphExists(ctx, "graph_demployment")
	if err != nil {
		log.Fatalf("error check graph 'graph_demployment', %v", err)
	}
	if !exists {
		_, err := GraphDB.CreateGraph(ctx, "graph_demployment", &driver.CreateGraphOptions{
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
	aql := `LET doc = ` + common.Stringfy(doc) + `
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

func RemoveProcInfoFromGraphDB(procLst []mutual.ProcessInfo) bool {
	ctx := context.Background()

	aql := `
FOR item in @procLst
FOR proc IN process
  FILTER proc._id == item.id

  LET tcp_set = (
	for tcp in conn_tcp
	  filter tcp._from == proc._id || tcp._to == proc._id
	  remove tcp in conn_tcp
	  return OLD._id
  )
  
  LET dep_set = (
    for dep in deployment
      filter dep._from == proc._id
      remove dep in deployment
      return OLD._id
  )
  
  remove proc in process
  
  return {
    tcp_set: tcp_set,
    dep_set: dep_set,
    proc: OLD._id
  }
`
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"procLst": procLst,
	})
	if err != nil {
		zap.L().Error("error remove process information from graphdb.", zap.Error(err))
		return false
	}
	defer cur.Close()

	return true
}

func SaveProcess(linuxIdentity string, procLst []mutual.ProcessInfo, timestamp int64) (map[string]string, error) {
	ctx := context.Background()
	aql := `
LET linuxId = @linuxId
LET timestamp = @timestamp
LET procLst = @procLst

FOR proc IN procLst
	LET doc = {
		"host_identity": linuxId,
		"pid":           proc.pid,
		"info": {
			"name":        proc.name,
			"ppid":        proc.ppid,
			"create_time": proc.create_time,
			"exec":        proc.exe,
			"cmd":         proc.cmd,
		},
		"timestamp": timestamp,
	}

	INSERT doc INTO process
	return {
		"id": NEW._id,
		"pid": NEW.pid,
		"identity": NEW.host_identity
	}
`

	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"procLst":   procLst,
		"linuxId":   linuxIdentity,
		"timestamp": timestamp,
	})
	if err != nil {
		zap.L().Error("error saving process lst", zap.Error(err))
		return nil, err
	}

	idMapping := make(map[string]string)
	for {
		info := map[string]any{
			"id":       "",
			"pid":      int32(-1),
			"identity": "",
		}
		_, err := cur.ReadDocument(context.Background(), &info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		idMapping[fmt.Sprintf("%s|%d", info["identity"].(string), int64(info["pid"].(float64)))] = info["id"].(string)
	}

	if cur != nil {
		cur.Close()
	}

	return idMapping, nil
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
	if cur != nil {
		cur.Close()
	}
	return nil
}

func SaveDeploymentRelation(linuxIdentity string, procKeyLst []string) {
	aql := `
LET linuxIdentity = @linuxIdentity
LET procLst = @procKeyLst

FOR h IN host
	FILTER h.host_identity == linuxIdentity
	for proc in procLst
		INSERT {
			"_from": proc, 
			"_to": h._id
		} INTO deployment
`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"linuxIdentity": linuxIdentity,
		"procKeyLst":    procKeyLst,
	})
	if err != nil {
		zap.L().Error("Failed to execute transaction: ", zap.Error(err))
	}
	defer cur.Close()
}

func SaveTCPConnection(relationLst []*TCPRelation) error {
	aql := `
FOR relation IN @relationLst
	FILTER relation.from != "" && relation.to != ""

	UPSERT {
		"_from": relation.from,
		"_to": relation.to
	} 
      INSERT {
		"_from": relation.from,
		"_to": relation.to,
		"local_identity": relation.local_identity,
		"local_pid": relation.local_pid,
		"local_ip": relation.local_ip,
		"local_port": relation.local_port,
		"remote_ip": relation.remote_ip,
		"remote_port": relation.remote_port,
		"remote_identity": relation.remote_identity,
		"remote_pid": relation.remote_pid,
		"timestamp": relation.timestamp
	  }
      UPDATE {
        "timestamp": relation.timestamp,
      }  
      IN conn_tcp
`
	ctx := context.Background()
	// zap.L().Debug("upsert connection relation: ", zap.String("input arg", common.ToString(linkLst)))
	cur, err := GraphDB.Query(ctx, aql, map[string]any{
		"relationLst": relationLst,
	})
	if err != nil {
		return err
	}
	defer cur.Close()
	return nil
}

func DeleteTimeoutTopo(timestamp int64) {
}

func QueryLinuxDesc(linuxId string) map[string]any {
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
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"id": linuxId,
	})
	if err != nil {
		log.Default().Panicln("error load linux description: ", err)
	}

	result := make(map[string]interface{})
	_, err = cur.ReadDocument(context.Background(), &result)
	if err != nil {
		log.Default().Panicln("error load linux description: ", err)
	}

	defer cur.Close()

	return result
}

func QueryLinuxTopo(linuxId string, showAll bool) ([]map[string]interface{}, error) {
	// options {"order": "bfs", uniqueEdges: "path"}
	aql := `
for h in host
  filter h.host_identity == @linuxId
  for v, e, p in 0..4 any h graph graph_demployment
	%s
  return {
    vertex: starts_with(v._id, "host/") ? {
          _id: v._id, 
          _key: v._key, 
          timestmap: v.timestamp, 
          name: v.name, 
          host_identity: v.host_identity,
          info: {
            hostname: v.info.hostname
          }
        } : v,
    edge: e
  } 
`
	if showAll {
		aql = fmt.Sprintf(aql, "")
	} else {
		aql = fmt.Sprintf(aql, `  FILTER (
    IS_SAME_COLLECTION('process', v)
    ? LENGTH(
        FOR edge IN conn_tcp
          FILTER edge._from == v._id OR edge._to == v._id
          LIMIT 1
          RETURN true
      ) > 0
    : true
  )`)
	}
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"linuxId": linuxId,
	})
	if err != nil {
		return nil, err
	}

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

	if cur != nil {
		cur.Close()
	}

	return result, nil
}

func QueryBizTopo(bizId int64, min, max int32, vset, rset []string) ([]map[string]interface{}, error) {
	aql := `
for b in business
  filter b.id == @bizId
  for v, e, p in @min..@max any b graph 'graph_demployment'
  return p
`
	minDepth := min - 1
	if minDepth < 0 {
		minDepth = 0
	}
	maxDepth := max - 1
	if maxDepth < 0 {
		maxDepth = 0
	}
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]interface{}{
		"bizId": bizId,
		"min":   minDepth,
		"max":   maxDepth,
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

func BatchGetNIC() []*CIDRRecord {
	aql := `
for h in host
  filter h.interface
  for i in h.interface
    for addr in i.addrs
		return {
			"cidr": addr.addr,
			"linuxId": h.host_identity
		}
`
	ctx := context.Background()
	cur, err := GraphDB.Query(ctx, aql, map[string]any{})
	if err != nil {
		zap.L().Error("error in BatchGetNIC: ", zap.Error(err))
	}
	defer cur.Close()

	lst := make([]*CIDRRecord, 0)
	for {
		info := new(CIDRRecord)
		_, err := cur.ReadDocument(ctx, info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Default().Println("error in BatchGetNIC: ", err)
		}
		lst = append(lst, info)
	}
	return lst

}

func SaveBiz(biz *Business) {
	ctx := context.Background()
	coll, err := GraphDB.Collection(ctx, "business")
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
	cur, err := GraphDB.Query(ctx, aql, map[string]any{
		"bizId":     linux.Biz.Id,
		"linuxId":   linux.LinuxId,
		"timestamp": linux.UpdateTimestamp,
	})
	if err != nil {
		log.Default().Printf("Failed to execute transaction: %v", err)
	}
	if cur != nil {
		cur.Close()
	}

}

func UpdateLinuxInGraphDB(origin, linux *Linux) {
	ctx := context.Background()
	aql := `FOR biz IN business
  FILTER biz.id == @bizId
  LIMIT 1
  FOR linux IN host
    FILTER linux.host_identity == @linuxId
    LIMIT 1
	UPDATE linux._key with {"name": @name, "host_identity": @host_identity} IN host
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

	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"bizId":         linux.Biz.Id,
		"linuxId":       origin.LinuxId,
		"host_identity": linux.LinuxId,
		"name":          linux.Hostname,
		"timestamp":     linux.UpdateTimestamp,
	})
	if err != nil {
		zap.L().Error("Failed to execute transaction: ", zap.Error(err))
	}
	if cursor != nil {
		cursor.Close()
	}
}

func DeleteLinuxInGraphDB(linux *Linux) {
	ctx := context.Background()
	aql := `let set = (for linux in host
  filter linux.host_identity == @linuxId
  for vertex, edge, path in  1..3 any linux._id graph graph_demployment
    return {
      "v": {
        "id": vertex._id,
        "key": vertex._key
      }, 
      "e": {
        "id": edge._id,
        "key": edge._key
      }
    }
)

let set1 = UNIQUE(
  for r in set
    filter starts_with(r.e.id, "conn_tcp/")
    return r.e.key
)
let result1 = (
  for r in set1
    remove {_key: r} in conn_tcp
    return OLD._id
)

let set2 = UNIQUE(for r in set
  filter starts_with(r.e.id, "deployment/")
    return r.e.key
)
let result2 = (
  for r in set2
    remove {_key: r} in deployment
    return OLD._id
)

let set3 = UNIQUE(for r in set
  filter starts_with(r.e.id, "res_consumption/")
    return r.e.key)
let result3 = (
  for r in set3
    remove {_key: r} in res_consumption
    return OLD._id
)

let set4 = UNIQUE(
for r in set
  filter starts_with(r.v.id, "process/")
    return r.v.key
)
let result4 = (
 for r in set4
   remove {_key: r} in process
)

let set5 = UNIQUE(
for r in set
  filter starts_with(r.v.id, "host/")
    return r.v.key
)
let result5 = (
 for r in set5
   remove {_key: r} in host
)
return 1`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"linuxId": linux.LinuxId,
	})
	if err != nil {
		zap.L().Error("Failed to execute delete linux in graph database: ", zap.Error(err))
	}
	defer cursor.Close()
}

func QueryProcessIdFromGraphDB(args []map[string]any) map[string]string {
	ctx := context.Background()
	aql := `
let lst = @lst
for item in lst
	for proc in process
		filter item.pid == proc.pid && item.identity == proc.host_identity
		return {
			"id": proc._id,
			"pid": proc.pid,
			"identity": proc.host_identity
		}
`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"lst": args,
	})
	if err != nil {
		zap.L().Error("error query process id by pid and identity", zap.Error(err))
	}
	defer cursor.Close()

	mapping := make(map[string]string)
	for {
		info := make(map[string]any)
		_, err := cursor.ReadDocument(ctx, &info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			zap.L().Error("error in transfor process id", zap.Error(err))
		}
		mapping[fmt.Sprintf("%s|%d", info["identity"].(string), int64(info["pid"].(float64)))] = info["id"].(string)
	}

	return mapping
}

func DeleteBizFromGraphDB(bizId int) {
	ctx := context.Background()
	aql := `
FOR biz IN business
	FILTER biz.id == @bizId
  LET cns_lst = (
    FOR cns in res_consumption
      filter cns._from == biz._id
      remove cns in res_consumption
      return OLD._id
  )
  remove biz in business
  return {
    "biz_id": OLD._id,
    "relation_lst": cns_lst
  }
`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"bizId": bizId,
	})
	if err != nil {
		zap.L().Error("error delete biz and res_consumption by biz id", zap.Error(err))
	}
	defer cursor.Close()
}

func GetProcessLstByLinuxId(linuxId string) []mutual.ProcessInfo {
	ctx := context.Background()
	aql := `
for proc in process
  filter proc.host_identity == @linuxId
  return proc
`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"linuxId": linuxId,
	})
	if err != nil {
		zap.L().Error("error delete biz and res_consumption by biz id", zap.Error(err))
	}

	procLst := make([]mutual.ProcessInfo, 0)
	for {
		procMap := make(map[string]any)
		_, err := cursor.ReadDocument(ctx, &procMap)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			zap.L().Error("error in transfor process id", zap.Error(err))
		}
		map0, ok := procMap["info"]
		if ok {
			infoMap := map0.(map[string]any)
			procLst = append(procLst, mutual.ProcessInfo{
				Id:         procMap["_id"].(string),
				Pid:        int32(procMap["pid"].(float64)),
				Name:       string(infoMap["name"].(string)),
				Ppid:       int32(infoMap["ppid"].(float64)),
				Exe:        string(infoMap["exec"].(string)),
				Cmd:        string(infoMap["cmd"].(string)),
				CreateTime: int64(infoMap["create_time"].(float64)),
			})
		}

	}

	defer cursor.Close()
	return procLst
}

func GetLinuxGraph(limit int) []map[string]interface{} {
	ctx := context.Background()
	aql := `
for h in host
  for v, e, p in 3 any h graph graph_demployment
  for e1 in p.edges
    filter is_same_collection("conn_tcp", e1)
    filter is_same_collection("host", v) && h._key != v._key
		limit @limit
        return {
          "from": {"name": h.name, "identity": h.host_identity, "key": h._key, "_id": h._id},
          "to": {"name": v.name, "identity": v.host_identity, "key": v._key, "_id": v._id},
        }
`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"limit": limit,
	})
	if err != nil {
		zap.L().Error("error get linux graph", zap.Error(err))
	}
	defer cursor.Close()
	result := make([]map[string]interface{}, 0)
	for {
		info := make(map[string]interface{})
		_, err := cursor.ReadDocument(ctx, &info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			zap.L().Error("error in reading linux graph", zap.Error(err))
		}
		result = append(result, info)
	}
	return result
}

func GetAllLinuxNode(limit int) []map[string]interface{} {
	// ctx := context.Background()
	// aql := `for h1 in host
	//   limit @limit
	// return {"name": h1.name, "identity": h1.host_identity, "key": h1._key, "id": h1._id}`
	// cursor, err := GraphDB.Query(ctx, aql, map[string]any{
	// 	"limit": limit,
	// })
	// if err != nil {
	// 	zap.L().Error("error get all linux", zap.Error(err))
	// }
	// defer cursor.Close()
	// result := make([]map[string]interface{}, 0)
	// for {
	// 	info := make(map[string]interface{})
	// 	_, err := cursor.ReadDocument(ctx, &info)
	// 	if driver.IsNoMoreDocuments(err) {
	// 		break
	// 	} else if err != nil {
	// 		zap.L().Error("error in reading all linux", zap.Error(err))
	// 	}
	// 	result = append(result, info)
	// }
	// return result
	return nil
}

func GetLinuxGraphWithStart(start string, depth int) []map[string]interface{} {
	ctx := context.Background()
	aql := `
  for v, e, p in 0..@depth any @start graph graph_demployment
		FILTER (
    IS_SAME_COLLECTION('process', v)
    ? LENGTH(
        FOR edge IN conn_tcp
          FILTER edge._from == v._id OR edge._to == v._id
          LIMIT 1
          RETURN true
      ) > 0
    : true)
  
        return {
          "vertex": is_same_collection("host", v) ? {"name": v.name, "identity": v.host_identity, "key": v._key, "_id": v._id} : v,
          "edge": e,
        }
`
	cursor, err := GraphDB.Query(ctx, aql, map[string]any{
		"start": start,
		"depth": depth,
	})
	if err != nil {
		zap.L().Error("error get linux graph", zap.Error(err))
	}
	defer cursor.Close()
	result := make([]map[string]interface{}, 0)
	for {
		info := make(map[string]interface{})
		_, err := cursor.ReadDocument(ctx, &info)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			zap.L().Error("error in reading linux graph", zap.Error(err))
		}
		result = append(result, info)
	}
	return result
}
