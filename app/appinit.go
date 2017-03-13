package app

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/go-xorm/xorm"

	"github.com/corvinusz/echo-xorm/logger"
	"github.com/corvinusz/echo-xorm/server/users"
)

// readConfig reads configuration file into application Config structure and inits in-memory token storage
func (a *Application) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return errors.New("Configuration file read error: " + cfgFileName + "\nError:" + err.Error())
	}
	_, err = toml.Decode(string(tomlData[:]), &a.C.Config)
	if err != nil {
		return errors.New("Configuration file decoding error: " + cfgFileName + "\nError:" + err.Error())
	}
	// init Logging data
	if len(a.C.Config.Logging.ID) == 0 {
		a.C.Config.Logging.ID = strconv.Itoa(os.Getpid())
	}
	if len(a.C.Config.Logging.LogTag) == 0 {
		a.C.Config.Logging.LogTag = os.Args[0]
	}
	return nil
}

// setupLogger sets apllication Logger up according to configuration settings
func (a *Application) initLogger() error {
	if a.C.Config.Logging.LogMode == "nil" || a.C.Config.Logging.LogMode == "null" {
		a.C.Logger = logger.NewNilLogger()
		return nil
	}
	a.C.Logger = logger.NewStdLogger(a.C.Config.Logging.ID, a.C.Config.Logging.LogTag)
	return nil
}

// init database
func (a *Application) initOrm() error {
	var err error
	// open database
	a.C.Orm, err = xorm.NewEngine(a.C.Config.Database.Db, a.C.Config.Database.Dsn)
	if err != nil {
		return err
	}
	// turn on logs
	a.C.Orm.ShowSQL(true)
	// migrate
	err = a.migrateDb()
	if err != nil {
		return err
	}
	// init data
	err = a.initDbData()
	return err
}

// migrate database
func (a *Application) migrateDb() error {
	var err error
	// migrate tables
	err = a.C.Orm.Sync(new(users.User))
	return err
}

// initDbData installs hardcoded data from config
func (a *Application) initDbData() error {
	user := &users.User{Login: "admin", Password: "admin"} // aaaa, backdoor
	_, err := user.Save(a.C.Orm)
	return err
}
