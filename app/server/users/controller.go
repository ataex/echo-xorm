package users

import (
	"net/http"
	"strconv"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/pkg/utils"

	"github.com/labstack/echo/v4"
)

// PostBody represents payload data format
type PostBody struct {
	Email       string  `json:"email"`
	DisplayName string  `json:"display_name"`
	PasswordURL *string `json:"password_url"`
	Password    string  `json:"password"`
}

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// GetAllUsers is a GET /users handler
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := FindAll(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser is a GET /users/{id} handler
func (h *Handler) GetUser(c echo.Context) error {
	var (
		user User
		err  error
	)

	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request paramer read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}

	err = user.FindOne(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser is a POST /users handler
func (h *Handler) CreateUser(c echo.Context) error {
	var body PostBody
	err := c.Bind(&body)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// validate
	if len(body.Email) == 0 {
		err = errors.NewWithCode(http.StatusBadRequest, "body validation error; email not recognized")
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	if len(body.Password) == 0 {
		err = errors.NewWithCode(http.StatusBadRequest, "body validation error; password not recognized")
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}

	// create
	user, err := NewUser(&body)
	if err != nil {
		return c.String(errors.Decompose(err))
	}
	// save
	err = user.Save(h.C.Orm)
	if err != nil {
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusCreated, user)
}

// PutUser is a PUT /users/{id} handler
func (h *Handler) PutUser(c echo.Context) error {
	var body PostBody
	// parse id
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request paramer read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// parse request body
	if err = c.Bind(&body); err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// construct user
	user := User{
		ID:       id,
		Email:    body.Email,
		Password: body.Password,
	}
	// update
	err = user.Update(h.C.Orm)
	if err != nil {
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser is a DELETE /users/{id} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var user User

	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request paramer read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}

	user.ID = id
	// delete
	err = user.Delete(h.C.Orm)
	if err != nil {
		return c.String(errors.Decompose(err))
	}
	return c.NoContent(http.StatusOK)
}
