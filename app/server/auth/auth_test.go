package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/corvinusz/echo-xorm/app/server/users"
	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func mockData() []*users.User {
	usrs := []*users.User{
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
	return usrs
}

func TestPostAuth(t *testing.T) {
	body := strings.NewReader(`{"email":"admin", "password":"admin"}`)
	rec, c, appc := unit.SetTestEnv(echo.POST, "/auth", body)
	h := NewHandler(appc)

	err := setDatabase(h)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPostAuthFail(t *testing.T) {
	body := strings.NewReader(`{"email":"admin", "password":"admin1"}`)
	rec, c, appc := unit.SetTestEnv(echo.POST, "/auth", body)
	h := NewHandler(appc)

	err := setDatabase(h)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func setDatabase(h *Handler) error {
	err := h.C.Orm.DropTables(&users.User{})
	if err != nil {
		return err
	}
	err = h.C.Orm.Sync(&users.User{})
	if err != nil {
		return err
	}
	_, err = h.C.Orm.Insert(mockData())
	if err != nil {
		return err
	}
	return nil
}
