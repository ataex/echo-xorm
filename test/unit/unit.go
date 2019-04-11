package unit

import (
	"fmt"
	"reflect"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/pkg/logger"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/objx"
)

func NewHandlerContext() *ctx.Context {
	testLogger := logger.NewStdLogger("echo-xorm-unit", "")
	testOrm, err := xorm.NewEngine("sqlite3", "/tmp/echo-xorm-test.sqlite")
	if err != nil {
		panic(err)
	}
	return &ctx.Context{
		Logger: testLogger,
		Orm:    testOrm,
	}
}

func ContainsJSON(s, substr string) (bool, error) {
	sMap, err := objx.FromJSON(s)
	if err != nil {
		return false, err
	}
	subMap, err := objx.FromJSON(substr)
	if err != nil {
		return false, err
	}

	for k := range subMap {
		if !reflect.DeepEqual(sMap[k], subMap[k]) {
			fmt.Printf("\nMISMATCHING KEY = %s\nMISMATCHING VALUE = %+v, %+v\n\n", k, sMap[k], subMap[k])
			return false, nil
		}
	}

	return true, nil
}
