package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
)

func main() {
	const cmpId = "gin"
	appCtx := appctx.NewAppContext(
		appctx.WithPrefix("gin"),
		appctx.WithName("Demo Gin"),
		appctx.WithComponent(ginserver.NewGinServer(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Error(err)
	}

	server := appCtx.MustGet(cmpId).(ginserver.GinServer)

	router := server.GetRouter()
	router.GET("/ping", demoHandler(appCtx))

	server.Start()
}

func demoHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := appCtx.Logger("demo")
		logger.Info("I am Iron Man")

		ctx.JSON(http.StatusOK, gin.H{"data": "pong"})
	}
}
