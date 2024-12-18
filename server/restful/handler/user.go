package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func login(username string, passwd string) (string, *model.User, error) {
	uLst := model.DBSelect("select id, username, is_active from user where username=? and passwd=?", username, passwd)
	uLen := len(uLst)
	if uLen > 1 {
		return "", nil, &common.InsightException{Code: 500, Msg: fmt.Sprintf("many user with username '%s'", username)}
	}
	if uLen < 1 {
		return "", nil, &common.InsightException{Code: 500, Msg: fmt.Sprintf("no user with username '%s'", username)}
	}
	item := uLst[0]
	if item["is_active"].(int64) != 1 {
		return "", nil, &common.InsightException{Code: 500, Msg: fmt.Sprintf("user with username '%s' is inactive", username)}
	}
	var user model.User
	user.ID = item["id"].(int64)
	user.Username = string(item["username"].([]uint8))
	user.IsActive = true
	token := uuid.NewString()
	bytes, _ := json.Marshal(user)
	model.CacheSet(token, string(bytes), common.SysArgs.Session.Expiration*time.Minute)
	return token, &user, nil
}

func UserLogin(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from: ", r)
			common.PrintStackTrace()
		}
	}()

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at user login.", zap.Error(err))
		return
	}
	var info = make(map[string]string)
	err = json.Unmarshal(body, &info)
	if err != nil {
		zap.L().Error("error parse json body at user login.", zap.Error(err))
		return
	}

	username := info["username"]
	passwd := info["passwd"]
	token, user, err := login(username, passwd)

	if err != nil {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusUnauthorized, Msg: err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
			"token": token,
			"user":  user,
		}, Msg: "success"})
	}
}
