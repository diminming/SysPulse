package restful

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"syspulse/model"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateDBRecord(r *model.Database) {
	sql := "insert into db_record(`name`, `db_id`, `type`, `biz_id`, `linux_id`, create_timestamp, update_timestamp) value(?,?,?,?,?,?,?)"
	id := model.DBInsert(sql, r.Name, r.DBID, r.Type, r.Biz.Id, r.Linux.Id, r.CreateTimestamp, r.UpdateTimestamp)
	r.Id = id
}

func DeleteDBRecord(rId int) {
	sql := "delete from db_record where id =? "
	model.DBDelete(sql, rId)
}

func UpdateDBRecord(r *model.Database) {
	sql := "update db_record set `name`=?, `db_id`=?, `type`=?, `biz_id`=?, `linux_id`=?, `update_timestamp`=? where `id` = ?"
	model.DBUpdate(sql, r.Name, r.DBID, r.Type, r.Biz.Id, r.Linux.Id, r.UpdateTimestamp, r.Id)
}

func GetDBRecordTotal() int {
	return model.GetDBTotal()
}

func GetDBRecordByPage(first int, pageSize int) []model.Database {
	sql := "SELECT \n" +
		"    db.id,\n" +
		"    db.name,\n" +
		"    db.db_id,\n" +
		"    db.type,\n" +
		"    b.id biz_id,\n" +
		"    b.biz_name,\n" +
		"    l.id linux_id,\n" +
		"    l.hostname,\n" +
		"    db.create_timestamp,\n" +
		"    db.update_timestamp\n" +
		"FROM\n" +
		"    (SELECT \n" +
		"        `id`,\n" +
		"            `name`,\n" +
		"            `db_id`,\n" +
		"            `type`,\n" +
		"            `biz_id`,\n" +
		"            `linux_id`,\n" +
		"            create_timestamp,\n" +
		"            update_timestamp\n" +
		"    FROM\n" +
		"        db_record\n" +
		"    LIMIT ? , ?) db\n" +
		"        LEFT JOIN\n" +
		"    (SELECT \n" +
		"        id, biz_name\n" +
		"    FROM\n" +
		"        biz) b ON db.biz_id = b.id\n" +
		"        LEFT JOIN\n" +
		"    (SELECT \n" +
		"        id, hostname\n" +
		"    FROM\n" +
		"        linux) l ON db.linux_id = l.id"
	lst := model.DBSelect(sql, first, pageSize)

	var result = []model.Database{}
	for _, o := range lst {
		var biz = model.Business{}
		var linux = model.Linux{}
		if o["biz_id"] != nil {
			biz = model.Business{Id: o["biz_id"].(int64), BizName: string(o["biz_name"].([]uint8))}
		}
		if o["linux_id"] != nil {
			linux = model.Linux{Id: o["linux_id"].(int64), Hostname: string(o["hostname"].([]uint8))}
		}
		database := model.Database{
			Id:              o["id"].(int64),
			Name:            string(o["name"].([]uint8)),
			DBID:            string(o["db_id"].([]uint8)),
			Type:            string(o["type"].([]uint8)),
			Linux:           linux,
			Biz:             biz,
			UpdateTimestamp: o["update_timestamp"].(int64),
			CreateTimestamp: o["create_timestamp"].(int64),
		}
		result = append(result, database)

	}

	return result
}

func (ws *WebServer) MappingHandler4Database() {
	ws.Get("/db/page", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()
		page, err := strconv.Atoi(values.Get("page"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		}
		pageSize, err := strconv.Atoi(values.Get("pageSize"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		}
		lst := GetDBRecordByPage(page, pageSize)
		total := GetDBRecordTotal()
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
			"lst":   lst,
			"total": total,
		}, Msg: "success"})
	})

	ws.Post("/db", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var db = model.Database{}
		err = json.Unmarshal(body, &db)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.CreateTimestamp = time.Now().Unix()
		db.UpdateTimestamp = time.Now().Unix()
		CreateDBRecord(&db)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: &db, Msg: "success"})
	})

	ws.Put("/db", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var db = model.Database{}
		err = json.Unmarshal(body, &db)
		if err != nil {
			fmt.Println(err)
			return
		}

		db.UpdateTimestamp = time.Now().Unix()
		UpdateDBRecord(&db)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: &db, Msg: "success"})
	})

	ws.Delete("/db", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()
		id, err := strconv.Atoi(values.Get("db_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		}
		DeleteDBRecord(id)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "success"})
	})
}
