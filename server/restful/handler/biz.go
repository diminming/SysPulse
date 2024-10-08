package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"

	"github.com/gin-gonic/gin"
)

func CreateBiz(biz *model.Business) {
	sql := "insert into biz(biz_name, biz_id, biz_desc, create_timestamp, update_timestamp) value(?, ?, ?, ?, ?)"
	id := model.DBInsert(sql, biz.BizName, biz.BizId, biz.BizDesc, biz.CreateTimestamp, biz.UpdateTimestamp)
	biz.Id = id
}

func UpdateBiz(biz *model.Business) {
	sql := "update biz set biz_name=?, biz_id=?, biz_desc=?, update_timestamp=? where id=?"
	affected, err := model.DBUpdate(sql, biz.BizName, biz.BizId, biz.BizDesc, biz.UpdateTimestamp, biz.Id)
	if err != nil {
		panic(err)
	}
	log.Printf("affected: %d", affected)
}

func GetBizByPage(page int, pageSize int) []model.Business {
	first := page * pageSize
	sql := "select id, biz_name, biz_id,biz_desc, create_timestamp, update_timestamp from biz order by update_timestamp desc, id desc limit ?, ?"
	lst := model.DBSelect(sql, first, pageSize)
	result := []model.Business{}
	for _, o := range lst {
		item := model.Business{Id: o["id"].(int64), BizName: string(o["biz_name"].([]uint8)), BizId: string(o["biz_id"].([]uint8)), BizDesc: string(o["biz_desc"].([]uint8)), CreateTimestamp: o["create_timestamp"].(int64), UpdateTimestamp: o["update_timestamp"].(int64)}
		result = append(result, item)
	}
	return result
}

func GetBizTotal() int64 {
	return model.GetBizTotal()
}

func DeleteBiz(bizId int) {
	sql := "delete from biz where id = ?"
	model.DBDelete(sql, bizId)
}

func NewBizRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var biz = model.Business{}
	err = json.Unmarshal(body, &biz)
	if err != nil {
		fmt.Println(err)
		return
	}
	biz.CreateTimestamp = time.Now().UnixMilli()
	biz.UpdateTimestamp = time.Now().UnixMilli()
	CreateBiz(&biz)
	model.SaveBiz(&biz)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: &biz, Msg: "success"})
}

func UpdateBizRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var biz = model.Business{}
	err = json.Unmarshal(body, &biz)
	if err != nil {
		fmt.Println(err)
		return
	}
	biz.UpdateTimestamp = time.Now().Unix()
	UpdateBiz(&biz)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: &biz, Msg: "success"})
}

func GetBizLstByPage(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
	}
	lst := GetBizByPage(page, pageSize)
	total := GetBizTotal()
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
		"lst":   lst,
		"total": total,
	}, Msg: "success"})
}

func DeleteBizRecord(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	bizId, err := strconv.Atoi(values.Get("biz_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
	}
	DeleteBiz(bizId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func GetBizCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: model.GetBizTotal(), Msg: "success"})
}
