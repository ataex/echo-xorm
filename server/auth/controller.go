package auth

import (
	"net/http"
	"time"

	"github.com/corvinusz/echo-xorm/ctx"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"github.com/corvinusz/echo-xorm/server/users"
)

// Handler represents handlers for '/auth'
type Handler struct {
	C   *ctx.Context
	Key []byte
}

// authInput represents payload data format
type authInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// authResponse represents payload response format
type authResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

type versionResponse struct {
	ServerTime uint64 `json:"server_time"`
	Version    string `json:"version"`
}

/*
// handler for /version
func (h *authHandler) getVersion(c echo.Context) error {
	vr := versionResponse{
		ServerTime: uint64(time.Now().Unix()),
		Version:    time.Now().String(),
	}
	return c.JSON(http.StatusOK, vr)
}

type authResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}
*/

// PostAuth is handler for /auth
func (h *Handler) PostAuth(c echo.Context) error {
	var (
		input authInput
		user  users.User
		err   error
		found bool
	)

	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// find user
	user = users.User{Login: input.Login}
	found, err = user.Find(h.C.Orm)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	if !found {
		return c.String(http.StatusForbidden, "user not found")
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

	resp := authResponse{
		Result: "OK",
		Token:  t,
	}
	return c.JSON(http.StatusOK, resp)
}

/*
// Set custom claims
claims := &jwtCustomClaims{
    "Jon Snow",
    true,
    jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
    },
}

// Create token with claims
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// Generate encoded token and send it as response.
t, err := token.SignedString([]byte("secret"))
if err != nil {
    return err
}
return c.JSON(http.StatusOK, echo.Map{
    "token": t,
})*/
