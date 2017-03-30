/*
 * Copyright (c) New Cloud Technologies, Ltd., 2013-2017
 *
 * You can not use the contents of the file in any way without New Cloud Technologies, Ltd. written permission.
 * To obtain such a permit, you should contact New Cloud Technologies, Ltd. at http://ncloudtech.com/contact.html
 *
 */

package main // license/lsx/cmd/lsx-admin-backend/main.go

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"github.com/corvinusz/echo-xorm/app"
	"github.com/corvinusz/echo-xorm/ctx"
)

var (
	configFlag = flag.String("config",
		"/usr/local/etc/echo-xorm-config.toml",
		"-config=\"path-to-your-config-file\" ")
)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	// parse flags
	flag.Parse()

	var (
		err error
		a   *app.Application
	)

	flags := &ctx.Flags{
		CfgFileName: *configFlag,
	}

	// create application
	a, err = app.New(flags)
	if err != nil {
		log.Fatal("error ", os.Args[0]+" initialization error: "+err.Error())
		os.Exit(1)
	}
	// setup OS-signal catchers
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() { // start OS-signal catching route
		for sig := range signalChannel {
			if a.C.Orm != nil {
				err = a.C.Orm.Close()
				if err != nil {
					a.C.Logger.Err("appcontrol", os.Args[0]+" db closing error on "+sig.String())
				}
			}
			if a.C.Logger != nil {
				a.C.Logger.Info("appcontrol", os.Args[0]+" graceful shutdown on "+sig.String())
				a.C.Logger.Close()
			}
			os.Exit(1)
		}
	}()

	// run application server
	if a.C.Logger == nil {
		log.Fatal("error ", os.Args[0]+" startup error: logger not initialized ")
		os.Exit(1)
	}
	a.C.Logger.Info("appcontrol", "started on localhost:"+a.C.Config.Port)
	a.Run()
}
