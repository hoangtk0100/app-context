package ginserver

import (
	"flag"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
)

const (
	defaultServerAddress = ":3000"
	defaultMode          = "debug"
)

type GinServer interface {
	GetAddress() string
	GetRouter() *gin.Engine
	Start()
}

type config struct {
	address string
	mode    string
}

type ginServer struct {
	id     string
	name   string
	router *gin.Engine
	logger appctx.Logger
	*config
}

func NewGinServer(id string) *ginServer {
	return &ginServer{
		id:     id,
		config: new(config),
	}
}

func (gs *ginServer) ID() string {
	return gs.id
}

func (gs *ginServer) InitFlags() {
	flag.StringVar(
		&gs.address,
		"gin-address",
		defaultServerAddress,
		"Gin server address - Default: 3000",
	)

	flag.StringVar(
		&gs.mode,
		"gin-mode",
		defaultMode,
		"Gin mode (debug | release) - Default: debug",
	)
}

func (gs *ginServer) Run(ac appctx.AppContext) error {
	gs.name = ac.GetName()
	gs.logger = ac.Logger(gs.id)

	if gs.mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	gs.router = gin.Default()
	gs.logger.Info("Init Gin server")

	return nil
}

func (gs *ginServer) Stop() error {
	return nil
}

func (gs *ginServer) GetAddress() string {
	return gs.address
}

func (gs *ginServer) GetRouter() *gin.Engine {
	return gs.router
}

func (gs *ginServer) Start() {
	gs.logger.Info("Start Gin server")

	if err := gs.router.Run(gs.address); err != nil {
		gs.logger.Fatal(err, "Cannot start server")
	}
}
