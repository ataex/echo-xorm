package users

import (
	"net/http"
	"testing"

	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func mockData() []*User {
	users := []*User{
		{
			ID:          1,
			Email:       "admin",
			DisplayName: "AdminName",
			Password:    "$2a$10$WUwK.b4F6BoXjBoq1ORpTONnXwrnoyA2EA7BfS9iNNEJRmkg8oGXq",
		},
		{
			ID:          100,
			Email:       "a_test_user_02@example.com",
			DisplayName: "a_test_user_02",
			Password:    "$2a$14$ZAolslKaP9AFy6PmxvZHQ.NIeZrMSQ0A/w65jpf4RRvTE4qyIvZ4C", // a_test_user_02
		},
	}
	return users
}

func setDatabase(h *Handler) error {
	err := h.C.Orm.DropTables(&User{})
	if err != nil {
		return err
	}
	err = h.C.Orm.Sync(&User{})
	if err != nil {
		return err
	}
	_, err = h.C.Orm.Insert(mockData())
	if err != nil {
		return err
	}
	return nil
}

func TestGetAll(t *testing.T) {
	rec, c, appc := unit.SetTestEnv(echo.GET, "/users", nil)
	h := NewHandler(appc)

	err := setDatabase(h)
	if err != nil {
		t.Fatal(err)
	}

	var expectedJSON = `[{"id":1,"email":"admin"}, {"id":100,"email":"a@a.a"}]`

	if assert.NoError(t, h.GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), expectedJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}

func TestGetUser(t *testing.T) {
	rec, c, appc := unit.SetTestEnv(echo.GET, "/users", nil)
	h := NewHandler(appc)

	err := setDatabase(h)
	if err != nil {
		t.Fatal(err)
	}

	c.SetParamNames("id")
	c.SetParamValues("100")

	var expectedJSON = `{"id":100,"email":"a_test_user_02@example.com", "displayName":"a_test_user_02"}`

	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), expectedJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}
