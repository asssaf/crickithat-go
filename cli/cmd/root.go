package cmd

import (
	"flag"

	"github.com/asssaf/crickithat-go/cli/cmd/servo"
	"github.com/asssaf/crickithat-go/cli/util"
)

type Command = util.Command

type RootCommand struct {
	*util.CompositeCommand
}

func NewRootCommand(usagePrefix string) *RootCommand {
	c := &RootCommand{
		CompositeCommand: util.NewCompositeCommand(
			flag.NewFlagSet(usagePrefix, flag.ExitOnError),
			[]Command{
				NewResetCommand(),
				servo.NewServoCommand(usagePrefix),
			},
			"",
		),
	}

	return c
}
