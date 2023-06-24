package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, err error) {
	if cErr, ok := err.(StatusCodeCarrier); ok {
		ctx.JSON(cErr.StatusCode(), cErr)
		return
	}

	ctx.JSON(
		http.StatusInternalServerError,
		ErrInternalServerError.WithError(err.Error()),
	)
}
