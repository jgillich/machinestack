package cmd

import (
	"flag"
	"fmt"

	"gitlab.com/faststack/machinestack/api"
	"gitlab.com/faststack/machinestack/config"
	"github.com/mitchellh/cli"
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

	a, err := api.NewServer(cfg)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if err := a.Start(); err != nil {
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
