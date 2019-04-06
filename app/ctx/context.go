package ctx

import (
	"github.com/go-xorm/xorm"

	"github.com/corvinusz/echo-xorm/pkg/logger"
)

// Context is a gate to application services
type Context struct {
	Orm        *xorm.Engine
	Logger     logger.Logger
	Config     *Config
	Flags      *Flags
	JWTSignKey []byte
}

// Flags represents start mode parameters for application
type Flags struct {
	CfgFileName string
}

// Config is a storage for admin application configuration
type Config struct {
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
