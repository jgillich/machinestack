package main

import (
	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/handler"
	"github.com/faststackco/machinestack/scheduler"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	Address   string
	TLSConfig *config.TLSConfig
	echo.Echo
}

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

	return &Server{config.Address, config.TLSConfig, *echo}, nil
}

// Start server while auto detecting TLS settings
func (s *Server) StartAuto() error {
	if s.TLSConfig.Enable {
		if s.TLSConfig.Auto {
			return s.StartAutoTLS(s.Address)
		}
		return s.StartTLS(s.Address, s.TLSConfig.Cert, s.TLSConfig.Key)
	}
	return s.Start(s.Address)
}
