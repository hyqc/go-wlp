package serve

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/middleware"
)

type Serve struct {
}

const HeartbeatPath = "/check"

func (s Serve) init() *gin.Engine {
	engine := gin.Default()
	//健康检查
	engine.GET(HeartbeatPath, check)
	auth := engine.Group("/debug/v1", middleware.GinJwt(global.Srv.User, global.Jwt), log.GinLog())

	return engine
}

func check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}

func (s Serve) Start() error {
	return nil
}
