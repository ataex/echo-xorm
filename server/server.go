package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server/auth"
	"github.com/corvinusz/echo-xorm/server/users"
	"github.com/corvinusz/echo-xorm/server/version"
)

// Server is an main application object that shared (read-only) to application modules
type Server struct {
	context    *ctx.Context
	signingKey []byte
}

// New constructor
func New(c *ctx.Context) *Server {
	s := new(Server)
	s.context = c
	s.signingKey = []byte(c.Config.Secret)
	return s
}

// Run registers API and starts http-server
func (s *Server) Run() {
	// Echo instance
	e := echo.New()

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var (
		authHandler    = auth.Handler{C: s.context, Key: s.signingKey}
		versionHandler = version.Handler{C: s.context}
		usersHandler   = users.Handler{C: s.context}
	)

	// Non-authored routes
	e.POST("/auth", authHandler.PostAuth)
	e.GET("/", versionHandler.GetVersion)
	e.GET("/version", versionHandler.GetVersion)
	// restricted
	r := e.Group("")
	// group middleware
	r.Use(middleware.JWT(s.signingKey))
	// users
	r.POST("/users", usersHandler.CreateUser)
	r.GET("/users", usersHandler.GetAllUsers)
	r.GET("/users/:id", usersHandler.GetUser)
	r.PUT("/users/:id", usersHandler.PutUser)
	r.DELETE("/users/:id", usersHandler.DeleteUser)

	// start server
	e.Server.Addr = ":" + s.context.Config.Port
	s.context.Logger.Info("appcontrol", "starting server at "+e.Server.Addr)
	err := e.Start(e.Server.Addr)
	if err != nil {
		s.context.Logger.Error("appcontrol", err.Error())
	}
}
