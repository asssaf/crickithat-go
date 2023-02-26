package servo

import (
	"flag"
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/crickithat-go/crickithat"
)

type StopCommand struct {
	fs    *flag.FlagSet
	num   int
	value int
}

func NewStopCommand() *StopCommand {
	c := &StopCommand{
		fs: flag.NewFlagSet("stop", flag.ExitOnError),
	}

	c.fs.IntVar(&c.num, "num", 0, "Servo number (1-4)")

	return c
}

func (c *StopCommand) Name() string {
	return c.fs.Name()
}

func (c *StopCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.num < 1 || c.num > 4 {
		return fmt.Errorf("servo num must be in the range 1-4: %d", c.num)
	}

	return nil
}

func (c *StopCommand) Execute() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatalf("host init: %w", err)
	}

	i2cPort, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		log.Fatalf("i2c open: %w", err)
	}

	opts := crickithat.DefaultOpts

	dev, err := crickithat.NewI2C(i2cPort, &opts)
	if err != nil {
		log.Fatal("device creation: %w", err)
	}
	defer dev.Halt()

	if err := dev.Init(); err != nil {
		log.Fatalf("device init: %w", err)
	}

	if err := dev.WriteServoStop(c.num - 1); err != nil {
		log.Fatalf("write servo: %w", err)
	}

	return nil
}
