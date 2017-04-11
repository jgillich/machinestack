package cmd

import (
	"flag"
	"fmt"

	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/model"
	"github.com/go-pg/migrations"
	"github.com/mitchellh/cli"
)

// MigrateCommand applies migrations
type MigrateCommand struct {
	Cli *cli.CLI
}

// Run nolint
func (c MigrateCommand) Run(args []string) int {
	var configPath = flag.String("config", "config.hcl", "config file path")

	cfg, err := config.ParseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	db, err := model.Db(cfg.PostgresConfig)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	_, _, err = migrations.Run(db, args...)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

// Help nolint
func (c MigrateCommand) Help() string {
	return c.Cli.HelpFunc(c.Cli.Commands) + "\n"
}

// Synopsis nolint
func (c MigrateCommand) Synopsis() string {
	return ""
}
