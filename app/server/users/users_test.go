package users

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/corvinusz/echo-xorm/test/unit"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func mockDbForGetAll(h *Handler) (*sql.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	rowsAll := sqlmock.NewRows([]string{"id", "email", "display_name"}).
		AddRow(1, "admin", "AdminName").
		AddRow(100, "a_test_user_02@example.com", "a_test_user_02")

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rowsAll)

	h.C.Orm.DB().DB = db
	return db, nil
}

func TestGetAll(t *testing.T) {
	rec, c, appc := unit.SetTestEnv(echo.GET, "/users", nil)
	h := NewHandler(appc)

	db, err := mockDbForGetAll(h)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var expectedJSON = `[{"id":1,"email":"admin"}, {"id":100,"email":"a@a.a"}]`

	if assert.NoError(t, h.GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), expectedJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}

func mockDbForGetOne(h *Handler) (*sql.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	rowsOne := sqlmock.NewRows([]string{"id", "email", "display_name"}).
		AddRow(100, "a_test_user_02@example.com", "a_test_user_02")

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rowsOne)

	h.C.Orm.DB().DB = db
	return db, nil
}

func TestGetOne(t *testing.T) {
	rec, c, appc := unit.SetTestEnv(echo.GET, "/users", nil)
	h := NewHandler(appc)
	c.SetParamNames("id")
	c.SetParamValues("100")

	db, err := mockDbForGetOne(h)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var expectedJSON = `{"id":100,"email":"a_test_user_02@example.com", "displayName":"a_test_user_02"}`

	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ok, err := unit.ContainsJSON(rec.Body.String(), expectedJSON)
		assert.Nil(t, err)
		assert.True(t, ok)
	}
}
