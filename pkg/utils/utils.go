package utils

import (
	"encoding/base64"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/pkg/errors"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/sha3"
)

// GetSHA3Hash return base64 encoded sha3hash
func GetSHA3Hash(data string) string {
	h := make([]byte, 64)
	sha3.ShakeSum256(h, []byte(data))
	return base64.StdEncoding.EncodeToString(h)
}

// GetEvent return string description of request from Echo.context
func GetEvent(c echo.Context) string {
	return c.Request().Method + " " + c.Path()
}

// BeginTansaction starts transaction for ORM with defined isolation LEVEL
// Returns transaction pointer and error (nil or prefixed)
// In case of error the transaction is always closed
func BeginTransaction(orm *xorm.Engine, level int) (*xorm.Session, error) {
	var query string
	tx := orm.NewSession()
	err := tx.Begin()
	if err != nil {
		tx.Close()
		return nil, errors.NewWithPrefix(err, "database error")
	}
	switch level {
	case ctx.LevelReadCommited:
		query = ";"
	case ctx.LevelRepeatableRead:
		query = "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"
	case ctx.LevelSerializable:
		query = "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE"
	}
	_, err = tx.Exec(query)
	if err != nil {
		tx.Close()
		return nil, errors.NewWithPrefix(err, "database error")
	}
	return tx, nil
}

// RollbackTransaction rolls transaction back
// returns error (nil or prefixed)
func RollbackTransaction(tx *xorm.Session, err error) error {
	erb := tx.Rollback()
	if erb != nil {
		if err != nil {
			erb = errors.NewWithPrefix(erb, err.Error())
		}
		return errors.NewWithPrefix(erb, "database error")
	}
	return errors.NewWithPrefix(err, "database error")
}
