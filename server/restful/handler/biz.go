package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"

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

func GetBizByPage(page int, pageSize int, keyword string) []model.Business {
	first := page * pageSize
	sqlstr := new(strings.Builder)
	sqlArgs := make([]any, 0, 10)
	sqlstr.WriteString("select id, biz_name, biz_id,biz_desc, create_timestamp, update_timestamp from biz\n")
	if keyword != "" && !(strings.TrimSpace(keyword) == "") {
		sqlstr.WriteString("where biz_name like ? or biz_id like ? or biz_desc like ?\n")
		likeArg := "%" + keyword + "%"
		sqlArgs = append(sqlArgs, likeArg)
		sqlArgs = append(sqlArgs, likeArg)
		sqlArgs = append(sqlArgs, likeArg)
	}
	sqlstr.WriteString("order by update_timestamp desc, id desc limit ?, ?")
	sqlArgs = append(sqlArgs, first)
	sqlArgs = append(sqlArgs, pageSize)
	lst := model.DBSelect(sqlstr.String(), sqlArgs...)
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

func DeleteBizFromDB(bizId int) {
	sql := "delete from biz where id = ?"
	model.DBDelete(sql, bizId)
}

func DeleteBiz(bizId int) {
	model.DeleteBizFromGraphDB(bizId)
	DeleteBizFromDB(bizId)
}

func NewBizRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error parsing request body: ", zap.Error(err))
		return
	}
	var biz = model.Business{}
	err = json.Unmarshal(body, &biz)
	if err != nil {
		zap.L().Error("error unpack request body: ", zap.Error(err))
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
	biz.UpdateTimestamp = time.Now().UnixMilli()
	UpdateBiz(&biz)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: &biz, Msg: "success"})
}

func CountInst(ctx *gin.Context) {
	bizId, err := strconv.ParseInt(ctx.Param("bizId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "bizId is not a number."})
	}
	count := model.CountInst(bizId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: count, Msg: "success"})
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

	keyword := values.Get("kw")
	zap.L().Info("biz page query: ", zap.String("keyword", keyword))

	lst := GetBizByPage(page, pageSize, keyword)
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

func QueryBizById(id int64) *model.Business {
	sqlstr := "select biz_name, biz_id, biz_desc from biz where id = ?"
	result := model.DBSelectRow(sqlstr, id)
	biz := new(model.Business)

	biz.Id = id
	biz.BizId = string(result["biz_id"].([]uint8))
	biz.BizDesc = string(result["biz_desc"].([]uint8))
	biz.BizName = string(result["biz_name"].([]uint8))

	return biz

}

func GetBizById(ctx *gin.Context) {
	bizId, err := strconv.ParseInt(ctx.Param("bizId"), 10, 64)
	if err != nil {
		zap.L().Error("can't process parameter biz_id from request: ", zap.Error(err))
	}
	biz := QueryBizById(int64(bizId))

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: biz, Msg: "success"})
}

func QueryBizTopo(ctx *gin.Context) {
	bizId, err := strconv.ParseInt(ctx.Param("bizId"), 10, 64)
	if err != nil {
		zap.L().Error("can't process parameter biz_id from request: ", zap.Error(err))
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error parsing request body: ", zap.Error(err))
		return
	}

	setting := struct {
		Min  int32    `json:"min"`
		Max  int32    `json:"max"`
		VSet []string `json:"vset"`
		ESet []string `json:"eset"`
	}{}
	err = json.Unmarshal(body, &setting)
	if err != nil {
		zap.L().Error("error unpack request body: ", zap.Error(err))
		return
	}

	result, err := model.QueryBizTopo(bizId, setting.Min, setting.Max, setting.VSet, setting.ESet)
	if err != nil {
		zap.L().Error("error query biz topo: ", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: result, Msg: "success"})
}
