package main

import (
	"log"
	"os"

	"github.com/faststackco/machinestack/cmd"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("api", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return cmd.RunCommand{Cli: c}, nil
		},
		"migrate": func() (cli.Command, error) {
			return cmd.MigrateCommand{Cli: c}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)

}
