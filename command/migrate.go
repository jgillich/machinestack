package cmd

import (
	"flag"
	"fmt"

	"gitlab.com/faststack/machinestack/config"
	"gitlab.com/faststack/machinestack/model"
	// Required for migrations to be picked up
	_ "gitlab.com/faststack/machinestack/model/migrations"
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

	old, new, err := migrations.Run(db, args...)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("Migrated: %v -> %v", old, new)

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
