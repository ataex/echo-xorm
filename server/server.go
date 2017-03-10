package server

import (
	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server/users"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server is an main application object that shared (read-only) to application modules
type Server struct {
	context *ctx.Context
}

// NewServer creates server
func NewServer(c *ctx.Context) *Server {
	s := new(Server)
	s.context = c
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
		usersHandler = users.Handler{C: s.context}
	)

	// Register routes
	// users
	e.POST("/users", usersHandler.CreateUser)
	e.GET("/users", usersHandler.GetAllUsers)
	e.GET("/users/:id", usersHandler.GetUser)
	e.PUT("/users/:id", usersHandler.PutUser)
	e.DELETE("/users/:id", usersHandler.DeleteUser)

	// start server
	s.context.Logger.Info("server started at localhost:" + s.context.Config.Port)
	s.context.Logger.Err(e.Start(":" + s.context.Config.Port))
}
