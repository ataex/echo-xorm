package version

import (
	"net/http"
	"testing"

	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	rec, c, appc := unit.SetTestEnv(echo.GET, "/version", nil)
	h := NewHandler(appc)

	var versionJSON = `{"version":"0.1.0develop", "result":"OK"}`

	if assert.NoError(t, h.GetVersion(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), versionJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}
