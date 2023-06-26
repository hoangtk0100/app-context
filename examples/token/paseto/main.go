package main

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
)

func main() {
	const cmpId = "paseto-token"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo PASETO Token"),
		appctx.WithComponent(token.NewPasetoMaker(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Fatal(err)
	}

	maker := appCtx.MustGet(cmpId).(core.TokenMakerComponent)

	accessToken, accessPayload, err := maker.CreateToken(token.AccessToken, "some-uid")
	if err != nil {
		log.Fatal(err)
	}

	verifiedPayload, err := maker.VerifyToken(accessToken)
	if err != nil {
		log.Error(err)
	}

	if accessPayload.UID != verifiedPayload.UID {
		log.Error(errors.New("Miss match UID"))
	}
	if !accessPayload.IssuedAt.Equal(verifiedPayload.IssuedAt) {
		log.Error(errors.New("Miss match IssuedAt"))
	}

	if !accessPayload.ExpiredAt.Equal(verifiedPayload.ExpiredAt) {
		log.Error(errors.New("Miss match ExpiredAt"))
	}
}
