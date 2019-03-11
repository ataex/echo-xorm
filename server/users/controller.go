package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/corvinusz/echo-xorm/ctx"
)

// PostBody represents payload data format
type PostBody struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// GetAllUsers is a GET /users handler
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := new(User).FindAll(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser is a GET /users/{id} handler
func (h *Handler) GetUser(c echo.Context) error {
	var (
		user   User
		err    error
		status int
	)

	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	status, err = user.Find(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser is a POST /users handler
func (h *Handler) CreateUser(c echo.Context) error {
	var (
		status int
		err    error
		user   User
		body   PostBody
	)

	if err = c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate
	if len(body.Email) == 0 {
		return c.String(http.StatusBadRequest, "login not recognized")
	}
	if len(body.Password) == 0 {
		return c.String(http.StatusBadRequest, "password not recognized")
	}

	// create
	user = User{
		Email:       body.Email,
		DisplayName: body.DisplayName,
		Password:    body.Password,
	}
	// save
	status, err = user.Save(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

// PutUser is a PUT /users/{id} handler
func (h *Handler) PutUser(c echo.Context) error {
	var (
		body   PostBody
		user   User
		id     uint64
		err    error
		status int
	)
	// parse id
	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// parse request body
	if err = c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// construct user
	user = User{
		ID:       id,
		Email:    body.Email,
		Password: body.Password,
	}
	// update
	status, err = user.Update(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser is a DELETE /users/{id} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var (
		id     uint64
		status int
		err    error
		user   User
	)

	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user.ID = id
	// delete
	status, err = user.Delete(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
