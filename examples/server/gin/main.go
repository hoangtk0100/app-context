package main

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
	"github.com/hoangtk0100/app-context/component/server/gin/middleware"
	"github.com/hoangtk0100/app-context/core"
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

	server := appCtx.MustGet(cmpId).(core.GinComponent)

	router := server.GetRouter()
	router.Use(middleware.Recovery(appCtx))
	router.GET("/ping", demoHandler(appCtx))
	router.GET("/error", demoErrorHandler(appCtx))

	server.Start()
}

func demoHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := appCtx.Logger("demo")
		logger.Info("I am Iron Man")

		core.SuccessResponse(ctx, core.NewDataResponse("pong"))
	}
}

func demoErrorHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		core.ErrorResponse(ctx, core.ErrBadRequest)
	}
}
