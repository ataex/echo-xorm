package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func mockData() []*User {
	users := []*User{
		{
			ID:    1,
			Email: "admin",
		},
		{
			ID:    100,
			Email: "a@.a.a",
		},
	}
	return users
}

func TestGetAll(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(unit.NewTestAppContext())

	c.SetPath("/users/")

	err := h.C.Orm.DropTables(&User{})
	if err != nil {
		t.Fatal(err)
	}
	err = h.C.Orm.Sync(&User{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = h.C.Orm.Insert(mockData())
	if err != nil {
		t.Fatal(err)
	}

	var userJSON = `[{"id":1,"email":"admin"}, {"id":100,"email":"a@a.a"}]`

	if assert.NoError(t, h.GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), userJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}

func TestGetUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(unit.NewTestAppContext())

	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("100")

	err := h.C.Orm.DropTables(&User{})
	if err != nil {
		t.Fatal(err)
	}
	err = h.C.Orm.Sync(&User{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = h.C.Orm.Insert(mockData())
	if err != nil {
		t.Fatal(err)
	}

	var userJSON = `{"id":100,"email":"a@.a.a"}`

	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), userJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}
