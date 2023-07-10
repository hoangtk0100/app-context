package grpcserver

import (
	"context"
	"strings"

	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
	authorizationCookie = "grpcgateway-cookie"
)

func AuthorizeUser(ctx context.Context, tokenMaker core.TokenMakerComponent) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, core.ErrUnauthorized.WithError(ErrMetadataMissing.Error())
	}

	// Check from Header
	values, ok := md[authorizationHeader]
	if ok {
		if len(values) == 0 {
			return nil, core.ErrUnauthorized.WithError(ErrAuthHeaderMissing.Error())
		}

		return getPayload(tokenMaker, values[0])
	}

	// Check from Cookie
	for _, cookies := range md[authorizationCookie] {
		for _, cookie := range strings.Split(cookies, ";") {
			fields := strings.Split(strings.TrimSpace(cookie), "=")
			if len(fields) == 2 && strings.ToLower(fields[0]) == authorizationHeader {
				return getPayload(tokenMaker, fields[1])
			}
		}
	}

	return nil, core.ErrUnauthorized.WithError(ErrAuthHeaderMissing.Error())
}

func getPayload(tokenMaker core.TokenMakerComponent, header string) (*token.Payload, error) {
	// <authorization-type><authorization-data>
	fields := strings.Fields(header)
	if len(fields) < 2 {
		return nil, core.ErrUnauthorized.WithError(ErrAuthHeaderInvalid.Error())
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, core.ErrUnauthorized.WithError(ErrAuthTypeUnsupported.Error())
	}

	accessToken := fields[1]
	payload, err := tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, core.ErrUnauthorized.WithError(ErrAccessTokenInvalid.Error())
	}

	return payload, nil
}
