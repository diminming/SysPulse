package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
)

func NewNMONRecord(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Default().Println(err)
		return
	}
	var nmon = new(model.NMON)
	err = json.Unmarshal(body, nmon)
	if err != nil {
		log.Default().Println(err)
		return
	}
	nmon.Source = "handler"
	nmon.CreateTimestamp = time.Now().UnixMilli()

	model.InsertNMONRecord(nmon)

	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: nmon, Msg: "success"})
}

func GetNMONRecordById(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	nmon := model.GetNMONRecordById(id)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: nmon, Msg: "success"})
}

func GetNMONRecordByPage(ctx *gin.Context) {
	id4Linux, _ := strconv.ParseInt(ctx.Query("linuxId"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	nmonLst := model.GetNMONRecordByPage(id4Linux, page, size)
	total := model.GetNMONTotal()
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]any{
		"lst":   nmonLst,
		"total": total,
	}, Msg: "success"})
}
