package main

import (
	"log"
	"os"

	"gitlab.com/faststack/machinestack/command"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("api", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return command.RunCommand{Cli: c}, nil
		},
		"migrate": func() (cli.Command, error) {
			return command.MigrateCommand{Cli: c}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)

}
