package servo

import (
	"errors"
	"flag"
	"fmt"
)

type Command interface {
	Init([]string) error
	Execute() error
	Name() string
}

type ServoCommand struct {
	fs       *flag.FlagSet
	commands []Command
	args     []string
}

func NewServoCommand() *ServoCommand {
	c := &ServoCommand{
		fs: flag.NewFlagSet("servo", flag.ExitOnError),
		commands: []Command{
			NewMoveCommand(),
		},
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: crickithat %s <command> ...\n", c.fs.Name())
		fmt.Fprintf(flag.CommandLine.Output(), "The commands are:\n")
		for _, c := range c.commands {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s\n", c.Name())
		}
		flag.PrintDefaults()
	}

	return c
}

func (c *ServoCommand) Name() string {
	return c.fs.Name()
}

func (c *ServoCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *ServoCommand) Execute() error {
	args := c.fs.Args()
	flag := c.fs

	command := flag.Arg(0)
	if command == "" {
		return errors.New("Missing command")
	}

	for _, c := range c.commands {
		if command == c.Name() {
			if err := c.Init(args[1:]); err != nil {
				return err
			}
			return c.Execute()
		}
	}

	return errors.New(fmt.Sprintf("unknown command: %s", command))
}
