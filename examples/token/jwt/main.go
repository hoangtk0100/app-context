package main

import (
	"time"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
)

func main() {
	const cmpId = "jwt-token"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo JWT Token"),
		appctx.WithComponent(token.NewJWTMaker(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Fatal(err)
	}

	maker := appCtx.MustGet(cmpId).(core.TokenMakerComponent)

	customToken, customPayload, err := maker.CreateToken(token.CustomToken, "some-uid", time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	verifiedPayload, err := maker.VerifyToken(customToken)
	if err != nil {
		log.Error(err)
	}

	if customPayload.UID != verifiedPayload.UID {
		log.Error(errors.New("Miss match UID"))
	}
	if !customPayload.IssuedAt.Equal(verifiedPayload.IssuedAt) {
		log.Error(errors.New("Miss match IssuedAt"))
	}

	if !customPayload.ExpiredAt.Equal(verifiedPayload.ExpiredAt) {
		log.Error(errors.New("Miss match ExpiredAt"))
	}
}
