package ctx

import (
	"github.com/corvinusz/echo-xorm/logger"
	"github.com/go-xorm/xorm"
)

// Context is a gate to application services
type Context struct {
	Orm    *xorm.Engine
	Logger logger.Logger
	Config *Config
	Flags  *Flags
}

// Flags represents start mode parameters for application
type Flags struct {
	CfgFileName string
}

// Config is a storage for admin application configuration
type Config struct {
	Secret   string `toml:"secret"`
	Version  string `toml:"version"`
	Port     string `toml:"port"`
	Database struct {
		Db  string `toml:"db"`
		Dsn string `toml:"dsn"`
	} `toml:"database"`
	Logging struct {
		LogMode string `toml:"log_mode"`
		LogTag  string `toml:"log_tag"`
		ID      string // will be process id
	} `toml:"logging"`
}
