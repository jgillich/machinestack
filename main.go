package main

import (
	"log"
	"os"

	"flag"
	"fmt"

	"github.com/faststackco/machinestack/config"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("api", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return RunCommand{cli: c}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)

}

// RunCommand nolint
type RunCommand struct {
	cli *cli.CLI
}

// Run nolint
func (c RunCommand) Run(args []string) int {
	var configPath = flag.String("config", "config.hcl", "config file path")

	cfg, err := config.ParseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	a, err := NewServer(cfg)
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
	return c.cli.HelpFunc(c.cli.Commands) + "\n"
}

// Synopsis nolint
func (c RunCommand) Synopsis() string {
	return ""
}
