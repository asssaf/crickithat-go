package util

import (
	"errors"
	"flag"
	"fmt"
)

type CompositeCommand struct {
	fs          *flag.FlagSet
	commands    []Command
	usagePrefix string
	args        []string
}

func NewCompositeCommand(fs *flag.FlagSet, subcommands []Command, usagePrefix string) *CompositeCommand {
	c := &CompositeCommand{
		fs:          fs,
		commands:    subcommands,
		usagePrefix: usagePrefix,
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s %s <command> ...\n", usagePrefix, c.fs.Name())
		fmt.Fprintf(flag.CommandLine.Output(), "The commands are:\n")
		for _, c := range c.commands {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s\n", c.Name())
		}
		flag.PrintDefaults()
	}

	return c
}

func (c *CompositeCommand) Name() string {
	return c.fs.Name()
}

func (c *CompositeCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *CompositeCommand) Execute() error {
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
