package ginserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/pflag"
)

const (
	defaultServerAddress = ":3000"
	defaultMode          = "debug"
)

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

func NewServer(id string) *ginServer {
	return &ginServer{
		id:     id,
		config: new(config),
	}
}

func (gs *ginServer) ID() string {
	return gs.id
}

func (gs *ginServer) InitFlags() {
	pflag.StringVar(
		&gs.address,
		"gin-address",
		defaultServerAddress,
		fmt.Sprintf("Gin server address - Default: %s", defaultServerAddress),
	)

	pflag.StringVar(
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

func (gs *ginServer) StartGracefully() {
	srv := &http.Server{
		Addr:    gs.address,
		Handler: gs.router,
	}

	go func() {
		gs.logger.Info("Server running at:", gs.address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			gs.logger.Fatal(err, "Server closed unexpectedly")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	gs.logger.Print("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		gs.logger.Fatal(err, "Server forced to shutdown")
	}

	gs.logger.Print("Server exited")
}
