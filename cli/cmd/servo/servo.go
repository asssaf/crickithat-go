package servo

import (
	"flag"
	"fmt"

	"github.com/asssaf/crickithat-go/cli/util"
)

type Command = util.Command

type ServoCommand struct {
	*util.CompositeCommand
}

func NewServoCommand(usagePrefix string) *ServoCommand {
	fs := flag.NewFlagSet("servo", flag.ExitOnError)
	c := &ServoCommand{
		CompositeCommand: util.NewCompositeCommand(
			fs,
			[]Command{
				NewMoveCommand(fmt.Sprintf("%s %s", usagePrefix, fs.Name())),
				NewStopCommand(),
			},
			usagePrefix,
		),
	}

	return c
}
