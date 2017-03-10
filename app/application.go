package app

import (
	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server"
)

// Application define a mode of running app
type Application struct {
	C *ctx.Context
}

// NewApplication constructor
func NewApplication(flags *ctx.Flags) (*Application, error) {
	app := new(Application)
	app.C = new(ctx.Context)
	// read config file
	err := app.initConfigFromFile(flags.CfgFileName)
	if err != nil {
		return nil, err
	}

	// init Logger
	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	// init Orm
	err = app.initOrm()
	return app, err
}

// Run starts application
func (a *Application) Run() {
	server := server.NewServer(a.C)
	server.Run()
}
