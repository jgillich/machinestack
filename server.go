package main

import (
	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/handler"
	"github.com/faststackco/machinestack/scheduler"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server is the http server that serves the api
type Server struct {
	Address   string
	TLSConfig *config.TLSConfig
	echo      *echo.Echo
}

// NewServer creates a new Server
func NewServer(config *config.Config) (*Server, error) {

	db := pg.Connect(&pg.Options{
		Addr:     config.PostgresConfig.Address,
		User:     config.PostgresConfig.Username,
		Password: config.PostgresConfig.Password,
		Database: config.PostgresConfig.Database,
	})

	sched, err := scheduler.NewScheduler(config.SchedulerConfig.Name, &config.DriverConfig.Options)
	if err != nil {
		return nil, err
	}

	hand := handler.NewHandler(db, sched)

	echo := echo.New()

	echo.Use(middleware.Gzip())

	echo.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  config.AuthConfig.Key,
		TokenLookup: "header:Authorization",
		Claims:      handler.JwtClaims{},
	}))

	echo.POST("/machines", hand.CreateMachine)
	echo.DELETE("/machines/:name", hand.DeleteMachine)

	echo.POST("/machines/:name/exec", hand.CreateExec)
	echo.GET("/exec/:id/io", hand.ExecIO)
	echo.GET("/exec/:id/control", hand.ExecControl)

	return &Server{config.Address, config.TLSConfig, echo}, nil
}

// Start starts the server
func (s *Server) Start() error {
	if s.TLSConfig != nil && s.TLSConfig.Enable {
		if s.TLSConfig.Auto {
			return s.echo.StartAutoTLS(s.Address)
		}
		return s.echo.StartTLS(s.Address, s.TLSConfig.Cert, s.TLSConfig.Key)
	}
	return s.echo.Start(s.Address)
}
