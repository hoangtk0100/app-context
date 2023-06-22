package main

import (
	appctx "github.com/hoangtk0100/app-context"
)

func main() {
	const cmpId = "abc"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo Component"),
		appctx.WithComponent(NewDemoComponent(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Error(err)
	}

	type CanDoSomething interface {
		GetData() string
		DoSomething() error
	}

	cmp := appCtx.MustGet(cmpId).(CanDoSomething)

	log.Print(cmp.GetData())
	_ = cmp.DoSomething()

	_ = appCtx.Stop()
}
