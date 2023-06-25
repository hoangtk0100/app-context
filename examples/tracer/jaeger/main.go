package main

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
	"github.com/hoangtk0100/app-context/component/server/gin/middleware"
	"github.com/hoangtk0100/app-context/component/tracer"
	"github.com/hoangtk0100/app-context/core"
	do "github.com/hoangtk0100/app-context/examples/tracer/jaeger/feature"
)

func main() {
	const ginId = "gin"
	const jaegerId = "jaeger"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo Jaeger"),
		appctx.WithComponent(tracer.NewJaeger(jaegerId)),
		appctx.WithComponent(ginserver.NewGinServer(ginId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Error(err)
	}

	server := appCtx.MustGet(ginId).(core.GinComponent)

	router := server.GetRouter()
	router.Use(middleware.Recovery(appCtx))
	router.GET("/ping", demoHandler(appCtx))

	server.Start()
}

func demoHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := appCtx.Logger("demo.do")
		do.DoSomething(ctx, logger)

		core.SuccessResponse(ctx, core.NewDataResponse("pong"))
	}
}
