package main

import (
	"fmt"
	"time"

	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/handler"
	"github.com/faststackco/machinestack/scheduler"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
		Addr:        config.PostgresConfig.Address,
		User:        config.PostgresConfig.Username,
		Password:    config.PostgresConfig.Password,
		Database:    config.PostgresConfig.Database,
		PoolSize:    20,
		PoolTimeout: time.Second * 5,
		ReadTimeout: time.Second * 5,
	})

	for _, model := range []interface{}{&handler.Machine{}} {
		if err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true, // TODO remove temp, use migrations
		}); err != nil {
			return nil, err
		}
	}

	fmt.Println(config.PostgresConfig)

	sched, err := scheduler.NewScheduler(config.SchedulerConfig.Name, &config.DriverConfig.Options)
	if err != nil {
		return nil, err
	}

	hand := handler.NewHandler(db, sched)

	e := echo.New()

	if config.LogLevel == "DEBUG" {
		e.Debug = true
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowOrigins,
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.Gzip())

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.JwtConfig.Secret),
		Claims:     &handler.JwtClaims{},
	}))

	e.GET("/machines", hand.MachineList)
	e.POST("/machines", hand.MachineCreate)
	e.GET("/machines/:name", hand.MachineInfo)
	e.DELETE("/machines/:name", hand.MachineDelete)

	e.POST("/machines/:name/exec", hand.ExecCreate)
	e.GET("/exec/:id/io", hand.ExecIO)
	e.GET("/exec/:id/control", hand.ExecControl)

	return &Server{config.Address, config.TLSConfig, e}, nil
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
