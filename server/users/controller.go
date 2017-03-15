package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/corvinusz/echo-xorm/ctx"
)

// UserInput represents payload data format
type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// GetAllUsers is a GET /users handler
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := new(User).FindAll(h.C.Orm)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser is a GET /users/{id} handler
func (h *Handler) GetUser(c echo.Context) error {
	var (
		user  User
		err   error
		found bool
	)

	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	found, err = user.Find(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	if !found {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser is a POST /users handler
func (h *Handler) CreateUser(c echo.Context) error {
	var (
		affected int64
		err      error
		user     User
		input    UserInput
	)

	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user = User{
		Login:    input.Login,
		Password: input.Password,
	}

	affected, err = user.Save(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	if affected == 0 {
		return c.String(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// PutUser is a PUT /users/{id} handler
func (h *Handler) PutUser(c echo.Context) error {
	var (
		input    UserInput
		user     User
		id       uint64
		err      error
		affected int64
	)
	// parse id
	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// parse request body
	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// construct user
	user = User{
		ID:       id,
		Login:    input.Login,
		Password: input.Password,
	}
	// update
	affected, err = user.Update(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	if affected == 0 {
		return c.String(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser is a DELETE /users/{id} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var (
		id       uint64
		affected int64
		err      error
		user     User
	)

	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user.ID = id
	// delete
	affected, err = user.Delete(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	if affected == 0 {
		return c.String(http.StatusConflict, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
