package server

import (
	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server/auth"
	"github.com/corvinusz/echo-xorm/server/users"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server is an main application object that shared (read-only) to application modules
type Server struct {
	context    *ctx.Context
	signingKey []byte
}

// NewServer constructor
func NewServer(c *ctx.Context) *Server {
	s := new(Server)
	s.context = c
	s.signingKey = []byte("secret")
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
		authHandler  = auth.Handler{C: s.context, Key: s.signingKey}
		usersHandler = users.Handler{C: s.context}
	)

	// Register routes
	e.POST("/auth", authHandler.PostAuth)
	// restricted
	r := e.Group("")
	r.Use(middleware.JWT(s.signingKey))
	// users
	r.POST("/users", usersHandler.CreateUser)
	r.GET("/users", usersHandler.GetAllUsers)
	r.GET("/users/:id", usersHandler.GetUser)
	r.PUT("/users/:id", usersHandler.PutUser)
	r.DELETE("/users/:id", usersHandler.DeleteUser)

	// start server
	s.context.Logger.Info("server started at localhost:" + s.context.Config.Port)
	s.context.Logger.Err(e.Start(":" + s.context.Config.Port))
}
