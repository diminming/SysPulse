package server

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/logging"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func NewRestfulServer() (*WebServer, error) {

	BuildAuthCache()

	router := gin.Default()
	router.Use(logging.GinLogger(), logging.GinRecovery(common.SysArgs.Logging.Level == "debug"))
	apiGrp := router.Group(common.SysArgs.Server.Restful.BasePath)
	callbackGrp := router.Group(common.SysArgs.Server.Restful.BasePathCallback)
	whiteLst := common.SysArgs.Server.Restful.WhiteLst

	apiGrp.Use(func(ctx *gin.Context) {
		path := ctx.FullPath()
		zap.L().Debug("request full path: ", zap.String("path", path))
		method := strings.ToLower(ctx.Request.Method)
		token := ctx.Request.Header.Get("token")

		if slices.Contains(whiteLst, method+":"+path) {
			ctx.Next()
			return
		} else if len(token) < 1 || !model.CacheExists(token) {
			ctx.JSON(http.StatusForbidden, response.JsonResponse{Status: http.StatusForbidden, Msg: "non-logging"})
			ctx.Abort()
			return
		}

		model.CacheExpire(token, common.SysArgs.Session.Expiration*time.Minute)
		userInfo := model.CacheGet(token)
		user := new(model.User)
		err := json.Unmarshal([]byte(userInfo), user)
		if err != nil {
			zap.L().Panic("error unmarshalling user info in auth method.", zap.Error(err))
		}

		if user.ID == 0 {
			ctx.Next()
			return
		} else {

			roleIdentityLst := make([]string, 0)
			for _, role := range user.RoleLst {
				roleIdentityLst = append(roleIdentityLst, role.Identity)
			}

			result := CheckAuth(path, method, roleIdentityLst)

			if result {
				ctx.Next()
				return
			}

			ctx.JSON(http.StatusUnauthorized, response.JsonResponse{Status: http.StatusUnauthorized, Msg: "no-auth"})
			ctx.Abort()
			return
		}

	})
	return &WebServer{
		router:      router,
		ApiGroup:    apiGrp,
		CallbackGrp: callbackGrp,
	}, nil
}

type WebServer struct {
	router      *gin.Engine
	ApiGroup    *gin.RouterGroup
	CallbackGrp *gin.RouterGroup
}

func (ws *WebServer) Post(url string, handler gin.HandlerFunc) {
	ws.ApiGroup.POST(url, handler)
}

func (ws *WebServer) Patch(url string, handler gin.HandlerFunc) {
	ws.ApiGroup.PATCH(url, handler)
}

func (ws *WebServer) PatchWithGrp(url string, handler gin.HandlerFunc, grp *gin.RouterGroup) {
	grp.PATCH(url, handler)
}

func (ws *WebServer) Put(url string, handler gin.HandlerFunc) {
	ws.ApiGroup.PUT(url, handler)
}

func (ws *WebServer) Get(url string, handler gin.HandlerFunc) {
	ws.ApiGroup.GET(url, handler)
}

func (ws *WebServer) Delete(url string, handler gin.HandlerFunc) {
	ws.ApiGroup.DELETE(url, handler)
}

func (ws *WebServer) Startup() {
	ws.SetupRoutes()
	ws.router.Run(common.SysArgs.Server.Restful.Addr)
}

func (ws *WebServer) Register4Api(method string, url string, handler gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		ws.ApiGroup.GET(url, handler)
	case http.MethodPost:
		ws.ApiGroup.POST(url, handler)
	case http.MethodPut:
		ws.ApiGroup.PUT(url, handler)
	case http.MethodPatch:
		ws.ApiGroup.PATCH(url, handler)
	case http.MethodDelete:
		ws.ApiGroup.DELETE(url, handler)
	}
}

func (ws *WebServer) Register4Callback(method string, url string, handler gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		ws.CallbackGrp.GET(url, handler)
	case http.MethodPost:
		ws.CallbackGrp.POST(url, handler)
	case http.MethodPut:
		ws.CallbackGrp.PUT(url, handler)
	case http.MethodPatch:
		ws.CallbackGrp.PATCH(url, handler)
	case http.MethodDelete:
		ws.CallbackGrp.DELETE(url, handler)
	}
}
