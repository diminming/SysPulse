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

func (ws *WebServer) MappingHandler4Cache() {
	ws.Get("/cache/page", func(ctx *gin.Context) {
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

	ws.Post("/cache", func(ctx *gin.Context) {
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

	ws.Put("/cache", func(ctx *gin.Context) {
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

	ws.Delete("/cache", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()
		id, err := strconv.Atoi(values.Get("cache_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		}
		DeleteDBRecord(id)
		ctx.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Msg: "success"})
	})
}
