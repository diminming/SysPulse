package restful

import (
	"fmt"
	"net/http"
	"syspulse/common"
	"syspulse/model"
	"time"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

type WebServer struct {
	router      *gin.Engine
	ApiGroup    *gin.RouterGroup
	CallbackGrp *gin.RouterGroup
}

func NewRestfulServer() (*WebServer, error) {
	router := gin.Default()
	apiGrp := router.Group(common.SysArgs.Server.Restful.BasePath)
	callbackGrp := router.Group(common.SysArgs.Server.Restful.BasePathCallback)

	apiGrp.Use(func(c *gin.Context) {
		path := c.FullPath()

		token := c.Request.Header.Get("token")

		if (len(token) < 1 || len(model.CacheGet(token)) < 1) && path != fmt.Sprintf("%s/login", common.SysArgs.Server.Restful.BasePath) {
			c.JSON(http.StatusForbidden, JsonResponse{Status: http.StatusForbidden, Msg: "non-logging"})
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

func (ws *WebServer) Mapping() {
	ws.MappingHandler4User()
	ws.MappingHandler4Biz()
	ws.MappingHandler4Linux()
	ws.MappingHandler4Database()
	ws.MappingHandler4Cache()
	ws.MappingRequest4Perfmance()
	ws.MappingHandler4Job()
}

func (ws *WebServer) Startup() {
	ws.Mapping()
	ws.router.Run(common.SysArgs.Server.Restful.Addr)
}
