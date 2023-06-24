package core

import appctx "github.com/hoangtk0100/app-context"

func Recovery() {
	if err := recover(); err != nil {
		appctx.GlobalLogger().GetLogger("recover").Error(err.(error), "Recovered")
	}
}
