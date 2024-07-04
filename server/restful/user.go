package restful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"syspulse/common"
	"syspulse/model"
	"time"

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

func (ws *WebServer) MappingHandler4User() {
	ws.Post("/login", func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from: ", r)
				common.PrintStackTrace()
			}
		}()

		username := c.Query("username")
		passwd := c.Query("passwd")
		token, user, err := login(username, passwd)

		if err != nil {
			c.JSON(http.StatusOK, JsonResponse{Status: http.StatusUnauthorized, Msg: err.Error()})
		} else {
			c.JSON(http.StatusOK, JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
				"token": token,
				"user":  user,
			}, Msg: "success"})
		}
	})
}
