package command

import (
	"flag"
	"fmt"

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

	db := cfg.PostgresConfig.Connect()

	sched, err := scheduler.NewScheduler(cfg.SchedulerConfig.Name, &cfg.DriverConfig.Options)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	handler := api.Handler{
		DB:           db,
		Scheduler:    sched,
		JWTSecret:    []byte(cfg.JwtConfig.Secret),
		AllowOrigins: cfg.AllowOrigins,
	}

	fmt.Printf("Serving on '%v'\n", cfg.Address)

	if err := handler.Serve(cfg.Address); err != nil {
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
