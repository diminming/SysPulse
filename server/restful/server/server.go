package server

import (
	"fmt"
	"net/http"
	"syspulse/common"
	"syspulse/model"
	"syspulse/restful/server/response"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRestfulServer() (*WebServer, error) {
	router := gin.Default()
	apiGrp := router.Group(common.SysArgs.Server.Restful.BasePath)
	callbackGrp := router.Group(common.SysArgs.Server.Restful.BasePathCallback)

	apiGrp.Use(func(c *gin.Context) {
		path := c.FullPath()

		token := c.Request.Header.Get("token")

		if (len(token) < 1 || len(model.CacheGet(token)) < 1) && path != fmt.Sprintf("%s/login", common.SysArgs.Server.Restful.BasePath) {
			c.JSON(http.StatusForbidden, response.JsonResponse{Status: http.StatusForbidden, Msg: "non-logging"})
			c.Abort()
			return
		} else {
			model.CacheExpire(token, common.SysArgs.Session.Expiration*time.Minute)
		}

		c.Next()
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
