package do

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"go.opencensus.io/trace"
)

func DoSomething(ctx context.Context, logger appctx.Logger) error {
	_, span := trace.StartSpan(ctx, "demo.do.something")
	defer span.End()

	logger.Info("Doing great")
	return nil
}
