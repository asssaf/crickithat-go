package servo

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/crickithat-go/crickithat"
)

type MoveCommand struct {
	fs     *flag.FlagSet
	num    int
	values []float64
}

func NewMoveCommand() *MoveCommand {
	c := &MoveCommand{
		fs: flag.NewFlagSet("move", flag.ExitOnError),
	}

	c.fs.IntVar(&c.num, "num", 0, "Servo number (1-4)")

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: crickithat servo %s <value> [<value> ... ]\n", c.fs.Name())
		fmt.Fprintf(flag.CommandLine.Output(), "Each <value> is in the range (0.0-1.0) as a factor of the pulse range.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "One or more <value>s can be provided and will be positioned in order.\n")
		c.fs.PrintDefaults()
	}

	return c
}

func (c *MoveCommand) Name() string {
	return c.fs.Name()
}

func (c *MoveCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.num < 1 || c.num > 4 {
		return fmt.Errorf("servo num must be in the range 1-4: %d", c.num)
	}

	if c.fs.NArg() == 0 {
		return fmt.Errorf("Missing value")
	}

	for _, arg := range c.fs.Args() {
		if value, err := strconv.ParseFloat(arg, 64); err != nil || value < 0.0 || value > 1.0 {
			return fmt.Errorf("servo value must be in the range 0.0-1.0: %s", arg)
		} else {
			c.values = append(c.values, value)
		}
	}

	return nil
}

func (c *MoveCommand) Execute() error {
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

	for _, value := range c.values {
		if err := dev.WriteServo(c.num-1, value); err != nil {
			log.Fatalf("write servo: %w", err)
		}
	}

	return nil
}
