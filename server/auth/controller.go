package auth

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server/users"
)

// Handler represents handlers for '/auth'
type Handler struct {
	C   *ctx.Context
	Key []byte
}

// Input represents payload data format
type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Result represents payload response format
type Result struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

// PostAuth is handler for /auth
func (h *Handler) PostAuth(c echo.Context) error {
	var (
		input Input
		user  users.User
		err   error
	)

	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// find user
	user = users.User{Login: input.Login}
	_, err = user.Find(h.C.Orm)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	//validate user credentials
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.String(http.StatusForbidden, "invalid credentials")
	}

	//create a HMAC SHA256 signer
	token := jwt.New(jwt.SigningMethodHS256)

	//set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "corvinusz/echo-xorm"
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).UTC().Unix()
	claims["aud"] = input.Login
	claims["jti"] = user.ID

	t, err := token.SignedString(h.Key)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, "Error while signing the token:"+err.Error())
	}

	resp := Result{
		Result: "OK",
		Token:  t,
	}
	return c.JSON(http.StatusOK, resp)
}
