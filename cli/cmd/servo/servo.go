package servo

import (
	"flag"

	"github.com/asssaf/crickithat-go/cli/util"
)

type Command = util.Command

type ServoCommand struct {
	*util.CompositeCommand
}

func NewServoCommand(usagePrefix string) *ServoCommand {
	c := &ServoCommand{
		CompositeCommand: util.NewCompositeCommand(
			flag.NewFlagSet("servo", flag.ExitOnError),
			[]Command{
				NewMoveCommand(),
				NewStopCommand(),
			},
			usagePrefix,
		),
	}

	return c
}
