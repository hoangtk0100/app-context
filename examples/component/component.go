package main

import (
	"github.com/spf13/pflag"

	appctx "github.com/hoangtk0100/app-context"
)

type demoComponent struct {
	id     string
	data   string
	logger appctx.Logger
}

func NewDemoComponent(id string) *demoComponent {
	return &demoComponent{id: id}
}

func (c *demoComponent) ID() string {
	return c.id
}

func (c *demoComponent) InitFlags() {
	pflag.StringVar(&c.data, "component-data", "demo", "Data string")
}

func (c *demoComponent) Run(ac appctx.AppContext) error {
	c.logger = ac.Logger(c.id)
	return nil
}

func (c *demoComponent) Stop() error {
	return nil
}

func (c *demoComponent) GetData() string {
	return c.data
}

func (c *demoComponent) DoSomething() error {
	c.logger.Print("LOL")
	return nil
}
