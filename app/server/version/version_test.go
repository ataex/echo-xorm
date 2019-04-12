package version

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(unit.NewTestAppContext())

	c.SetPath("/version")

	var versionJSON = `{"version":"0.1.0develop", "result":"OK"}`

	if assert.NoError(t, h.GetVersion(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), versionJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}
