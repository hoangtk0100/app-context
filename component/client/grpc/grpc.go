package grpcclient

import (
	"context"
	"fmt"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultClientAddress = "localhost:8088"
)

type config struct {
	address     string
	tlsCertFile string
	dialOpts    []grpc.DialOption
}

type grpcClient struct {
	id     string
	prefix string
	logger appctx.Logger
	*config
}

func NewClient(id string, prefix string) *grpcClient {
	return &grpcClient{
		id:     id,
		prefix: prefix,
		config: new(config),
	}
}

func (gc *grpcClient) ID() string {
	return gc.id
}

func (gc *grpcClient) InitFlags() {
	prefix := gc.prefix
	if prefix != "" {
		prefix = "-" + prefix
	}

	pflag.StringVar(
		&gc.address,
		fmt.Sprintf("grpc%s-client-address", prefix),
		defaultClientAddress,
		fmt.Sprintf("GRPC client%s address - Default: %s", prefix, defaultClientAddress),
	)

	pflag.StringVar(
		&gc.tlsCertFile,
		fmt.Sprintf("grpc%s-client-tls-cert-file", prefix),
		"",
		fmt.Sprintf("GRPC client%s TLS cert file", prefix),
	)
}

func (gc *grpcClient) Run(ac appctx.AppContext) error {
	gc.logger = ac.Logger(gc.id)

	if gc.tlsCertFile != "" {
		creds, err := credentials.NewClientTLSFromFile(gc.tlsCertFile, "")
		if err != nil {
			return core.ErrInternalServerError.
				WithError(ErrCannotAddClientTLS.Error()).
				WithDebug(err.Error())
		}

		gc.dialOpts = append(gc.dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		gc.logger.Warn("server: insecure mode.")
		gc.dialOpts = append(gc.dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	gc.dialOpts = append(gc.dialOpts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)))
	gc.logger.Infof("Init GRPC%s client", gc.getPrefixedValue())
	return nil
}

func (gc *grpcClient) Stop() error {
	return nil
}

func (gc *grpcClient) WithPrefix(prefix string) {
	gc.prefix = prefix
}

func (gc *grpcClient) getPrefixedValue() string {
	prefix := gc.prefix
	if prefix != "" {
		prefix = " " + prefix
	}

	return prefix
}

func (gc *grpcClient) WithAddress(address string) {
	gc.logger.Infof("GRPC%s client address is set to %s", gc.getPrefixedValue(), address)
	gc.address = address
}

func (gc *grpcClient) GetAddress() string {
	return gc.address
}

func (gc *grpcClient) GetLogger() appctx.Logger {
	return gc.logger
}

func (gc *grpcClient) Dial(options ...grpc.DialOption) *grpc.ClientConn {
	return gc.DialContext(context.Background(), options...)
}

func (gc *grpcClient) DialContext(ctx context.Context, options ...grpc.DialOption) *grpc.ClientConn {
	if len(options) > 0 {
		gc.dialOpts = append(gc.dialOpts, options...)
	}

	prefix := gc.getPrefixedValue()
	gc.logger.Infof("GRPC%s client dialing to address: %s", prefix, gc.address)

	conn, err := grpc.DialContext(ctx, gc.address, gc.dialOpts...)
	if err != nil {
		gc.logger.Errorf(err, "GRPC%s client dial to %s failed", prefix, gc.address)
		return nil
	}

	return conn
}
