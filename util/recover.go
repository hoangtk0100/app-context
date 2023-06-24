package util

import appctx "github.com/hoangtk0100/app-context"

func Recovery() {
	if err := recover(); err != nil {
		fErr := err.(error)
		appctx.GlobalLogger().GetLogger("recovery").Error(fErr, "Recovered: ", fErr.Error())
	}
}
