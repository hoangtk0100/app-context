package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultServerAddress   = "localhost:8088"
	defaultAPIPrefix       = "/"
	defaultSwaggerPrefix   = "/swagger/"
	defaultShutdownTimeout = time.Second * 10
)

type config struct {
	lis                        net.Listener
	address                    string
	tlsCertFile                string
	tlsKeyFile                 string
	apiPrefix                  string
	enableSwagger              bool
	swaggerPrefix              string
	enableMetrics              bool
	mapProtoResponseFieldStyle bool
	shutdownTimeout            time.Duration
	serverOptions              []grpc.ServerOption
	serveMuxOptions            []runtime.ServeMuxOption

	// Interceptors
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor

	dialOpts []grpc.DialOption
}

type grpcServer struct {
	id      string
	ac      appctx.AppContext
	logger  appctx.Logger
	server  *grpc.Server
	gateway *runtime.ServeMux

	*config
}

func NewServer(id string) *grpcServer {
	return &grpcServer{
		id:     id,
		config: new(config),
	}
}

func (gs *grpcServer) ID() string {
	return gs.id
}

func (gs *grpcServer) InitFlags() {
	pflag.StringVar(
		&gs.address,
		"grpc-server-address",
		defaultServerAddress,
		fmt.Sprintf("GRPC server address - Default: %q", defaultServerAddress),
	)

	pflag.StringVar(
		&gs.tlsCertFile,
		"grpc-server-tls-cert-file",
		"",
		"GRPC server TLS cert file",
	)

	pflag.StringVar(
		&gs.tlsKeyFile,
		"grpc-server-tls-key-file",
		"",
		"GRPC server TLS key file",
	)

	pflag.StringVar(
		&gs.apiPrefix,
		"grpc-server-api-prefix",
		defaultAPIPrefix,
		fmt.Sprintf("GRPC server API prefix - Default: %q", defaultAPIPrefix),
	)

	pflag.BoolVar(
		&gs.enableSwagger,
		"grpc-server-enable-swagger",
		false,
		"GRPC server enable Swagger - Default: false",
	)

	pflag.StringVar(
		&gs.swaggerPrefix,
		"grpc-server-swagger-prefix",
		defaultSwaggerPrefix,
		fmt.Sprintf("GRPC server Swagger prefix - Default: %q", defaultSwaggerPrefix),
	)

	pflag.BoolVar(
		&gs.enableMetrics,
		"grpc-server-enable-metrics",
		false,
		"GRPC server enable metrics (Prometheus) - Default: false",
	)

	pflag.BoolVar(
		&gs.mapProtoResponseFieldStyle,
		"grpc-server-map-proto-response-field-style",
		true,
		"GRPC server map proto response field style - Default: true (false: camelCase)",
	)

	pflag.DurationVar(
		&gs.shutdownTimeout,
		"grpc-server-shutdown-timeout",
		defaultShutdownTimeout,
		"GRPC server shutdown timeout - Default: 10s",
	)
}

func (gs *grpcServer) Run(ac appctx.AppContext) error {
	gs.ac = ac
	gs.logger = ac.Logger(gs.id)
	gs.lis = gs.getListener()

	if (gs.tlsCertFile != "" && gs.tlsKeyFile == "") ||
		(gs.tlsKeyFile != "" && gs.tlsCertFile == "") {
		gs.logger.Fatal(ErrTLSCertNotFull)
	}

	if gs.enableSwagger && gs.swaggerPrefix == "" {
		gs.logger.Fatal(ErrSwaggerPrefixMissing)
	}

	if gs.enableMetrics {
		gs.streamInterceptors = append(gs.streamInterceptors, grpc_prometheus.StreamServerInterceptor)
		gs.unaryInterceptors = append(gs.unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
	}

	gs.serverOptions = append(gs.serverOptions, grpc.UnaryInterceptor(GrpcLogger))

	if gs.isSecured() {
		creds, err := credentials.NewServerTLSFromFile(gs.tlsCertFile, gs.tlsKeyFile)
		if err != nil {
			gs.logger.Error(err, ErrCannotReadTLSCert.Error())
		}

		gs.serverOptions = append(gs.serverOptions, grpc.Creds(creds))
	}

	if gs.mapProtoResponseFieldStyle {
		// To map response field names the same style in proto files (default: camelCase)
		jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

		gs.serveMuxOptions = append(gs.serveMuxOptions, jsonOption)
	}

	if gs.isSecured() {
		creds, err := credentials.NewClientTLSFromFile(gs.tlsCertFile, "")
		if err != nil {
			return core.ErrInternalServerError.
				WithError(ErrCannotAddClientTLS.Error()).
				WithDebug(err.Error())
		}

		gs.dialOpts = append(gs.dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		gs.logger.Warn("server: insecure mode.")
		gs.dialOpts = append(gs.dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	gs.logger.Info("Init GRPC server")

	return nil
}

func (gs *grpcServer) Stop() error {
	return nil
}

func (gs *grpcServer) WithAddress(address string) {
	gs.logger.Infof("GRPC server listener is set to %s", address)
	gs.address = address
}

func (gs *grpcServer) WithListener(lis net.Listener) {
	newAddress := lis.Addr().String()
	if gs.address != "" {
		gs.logger.Infof("GRPC server listener is set to %s", newAddress)
		gs.address = newAddress
	}

	gs.lis = lis
}

func (gs *grpcServer) WithServerOptions(serverOpts ...grpc.ServerOption) {
	gs.serverOptions = append(gs.serverOptions, serverOpts...)
}

func (gs *grpcServer) WithServeMuxOptions(muxOpts ...runtime.ServeMuxOption) {
	gs.serveMuxOptions = append(gs.serveMuxOptions, muxOpts...)
}

func (gs *grpcServer) WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	gs.unaryInterceptors = append(gs.unaryInterceptors, interceptors...)
}

func (gs *grpcServer) WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	gs.streamInterceptors = append(gs.streamInterceptors, interceptors...)
}

func (gs *grpcServer) GetLogger() appctx.Logger {
	return gs.logger
}

func (gs *grpcServer) GetServer() *grpc.Server {
	gs.updateServerOptions()

	if gs.server == nil {
		gs.server = grpc.NewServer(gs.serverOptions...)
	}

	return gs.server
}

func (gs *grpcServer) GetGateway() *runtime.ServeMux {
	if gs.gateway == nil {
		gs.gateway = runtime.NewServeMux(gs.serveMuxOptions...)
	}

	return gs.gateway
}

func (gs *grpcServer) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	if gs.gateway != nil {
		gs.startGateway(ctx)
	} else {
		gs.startServer(ctx)
	}
}

func (gs *grpcServer) isSecured() bool {
	return gs.tlsCertFile != "" && gs.tlsKeyFile != ""
}

func (gs *grpcServer) getListener() net.Listener {
	if gs.lis != nil {
		return gs.lis
	}

	listener, err := net.Listen("tcp", gs.address)
	if err != nil {
		gs.logger.Fatal(err, ErrCannotCreateListener)
	}

	return listener
}

func (gs *grpcServer) serveSwagger(mux *http.ServeMux) {
	// Serve swagger from statik binary file
	// New() use default namespace
	// NewWithNamespace() for using custom namespace
	if gs.enableSwagger {
		statikFS, err := fs.New()
		if err != nil {
			gs.logger.Fatal(err, ErrCannotCreateStatikFS.Error())
		}

		swaggerHandler := http.StripPrefix(gs.swaggerPrefix, http.FileServer(statikFS))
		mux.Handle(gs.swaggerPrefix, swaggerHandler)
	}
}

func (gs *grpcServer) updateServerOptions() {
	if len(gs.streamInterceptors) > 0 {
		gs.serverOptions = append(gs.serverOptions, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(gs.streamInterceptors...)))
	}

	if len(gs.unaryInterceptors) > 0 {
		gs.serverOptions = append(gs.serverOptions, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(gs.unaryInterceptors...)))
	}
}

func (gs *grpcServer) startServer(ctx context.Context) {
	if gs.enableMetrics {
		grpc_prometheus.Register(gs.server)
	}

	gs.logger.Infof("Start GRPC server at %s", gs.address)

	reflection.Register(gs.server)
	err := gs.server.Serve(gs.lis)
	if err != nil {
		gs.logger.Fatal(err, ErrCannotStartServer.Error())
	}
}

func (gs *grpcServer) startGateway(ctx context.Context) {
	// Only executed before exiting this runGatewayServer function
	// Cancelling a context is a way to prevent the system from doing unnecessary works
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	// Convert HTTP request to GRPC format, reroute them to the GRPC mux
	mux.Handle(gs.apiPrefix, gs.gateway)

	gs.serveSwagger(mux)

	gs.logger.Infof("Start HTTP gateway server at %s", gs.address)

	err := http.Serve(gs.lis, HttpLogger(mux))
	if err != nil {
		gs.logger.Fatal(err, ErrCannotStartGatewayServer.Error())
	}
}
