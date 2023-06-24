package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	core "github.com/hoangtk0100/app-context/core"
)

func Recovery(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Header("Content-Type", "application/json")

				if appErr, ok := err.(core.StatusCodeCarrier); ok {
					ctx.AbortWithStatusJSON(appErr.StatusCode(), appErr)
				} else {
					// General panic cases
					ctx.AbortWithStatusJSON(
						http.StatusInternalServerError,
						core.ErrInternalServerError,
					)
				}

				appCtx.Logger("service").Errorf(err.(error), "%+v\n", err)

				if gin.IsDebugging() {
					panic(err)
				}
			}
		}()

		ctx.Next()
	}
}
