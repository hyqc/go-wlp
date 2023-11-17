package serve

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/app/account/controller"
)

type Serve struct {
}

const HeartbeatPath = "/check"

func (s Serve) init() *gin.Engine {
	engine := gin.Default()
	//健康检查
	engine.GET(HeartbeatPath, check)
	root := engine.Group("/")
	account := root.Group("/account")
	{
		account.POST("/login", controller.Login)
	}

	//auth := root.Group("/auth", middleware.JWTCheck(config.JWT))
	//{
	//
	//}

	return engine
}

func check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}

func (s Serve) Start() error {

	return nil
}
