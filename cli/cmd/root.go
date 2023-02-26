package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/asssaf/crickithat-go/cli/cmd/servo"
	"github.com/asssaf/crickithat-go/cli/util"
)

type Command = util.Command

func Execute() error {
	commands := []Command{
		NewResetCommand(),
		servo.NewServoCommand(),
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: crickithat <command> ...\n")
		fmt.Fprintf(flag.CommandLine.Output(), "The commands are:\n")
		for _, c := range commands {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s\n", c.Name())
		}
		flag.PrintDefaults()
	}

	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		return errors.New("Missing command")
	}

	args := flag.Args()
	for _, c := range commands {
		if command == c.Name() {
			if err := c.Init(args[1:]); err != nil {
				return err
			}
			return c.Execute()
		}
	}

	return errors.New(fmt.Sprintf("unknown command: %s", command))
}
