package cmd

import (
	"flag"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/mitchellh/cli"
	"gitlab.com/faststack/machinestack/api"
	"gitlab.com/faststack/machinestack/config"
	"gitlab.com/faststack/machinestack/scheduler"
)

// RunCommand is the default command that runs the server
type RunCommand struct {
	Cli *cli.CLI
}

// Run nolint
func (c RunCommand) Run(args []string) int {
	var configPath = flag.String("config", "config.hcl", "config file path")

	cfg, err := config.ParseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	db := pg.Connect(&pg.Options{
		Addr:     config.PostgresConfig.Address,
		User:     config.PostgresConfig.Username,
		Password: config.PostgresConfig.Password,
		Database: config.PostgresConfig.Database,
		//PoolSize:    20,
		//PoolTimeout: time.Second * 5,
		//ReadTimeout: time.Second * 5,
	})

	sched, err := scheduler.NewScheduler(config.SchedulerConfig.Name, &config.DriverConfig.Options)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	handler := api.Handler{
		DB:           db,
		Scheduler:    sched,
		JWTSecret:    cfg.JwtConfig.Secret,
		AllowOrigins: cfg.AllowOrigins,
	}

	if err := handler.Serve(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

// Help nolint
func (c RunCommand) Help() string {
	return c.Cli.HelpFunc(c.Cli.Commands) + "\n"
}

// Synopsis nolint
func (c RunCommand) Synopsis() string {
	return ""
}
