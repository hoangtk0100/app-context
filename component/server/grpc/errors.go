package grpcserver

import "errors"

var (
	ErrMetadataMissing          = errors.New("missing metadata")
	ErrAuthHeaderMissing        = errors.New("missing authorization header")
	ErrAuthHeaderInvalid        = errors.New("invalid authorization header format")
	ErrAuthTypeUnsupported      = errors.New("unsupported authorization type")
	ErrAccessTokenInvalid       = errors.New("invalid access token")
	ErrTLSCertNotFull           = errors.New("TLS cert or key file is missing")
	ErrCannotReadTLSCert        = errors.New("cannot read TLS cert or key file")
	ErrSwaggerPrefixMissing     = errors.New("missing swagger prefix")
	ErrCannotCreateListener     = errors.New("cannot create listener")
	ErrCannotStartServer        = errors.New("cannot start GRPC server")
	ErrCannotStartGatewayServer = errors.New("cannot start HTTP gateway server")
	ErrCannotCreateStatikFS     = errors.New("cannot create statik fs")
	ErrCannotAddClientTLS       = errors.New("cannot add client TLS")
)
