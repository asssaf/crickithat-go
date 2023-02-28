package cmd

import (
	"flag"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/crickithat-go/crickithat"
)

type NeopixelCommand struct {
	fs    *flag.FlagSet
	pin   int
	red   int
	green int
	blue  int
}

func NewNeopixelCommand() *NeopixelCommand {
	c := &NeopixelCommand{
		fs: flag.NewFlagSet("neopixel", flag.ExitOnError),
	}

	c.fs.IntVar(&c.pin, "pin", 27, "Pin")
	c.fs.IntVar(&c.red, "red", 0, "Red (0-255)")
	c.fs.IntVar(&c.green, "green", 0, "Green (0-255)")
	c.fs.IntVar(&c.blue, "blue", 0, "Blue (0-255)")

	return c
}

func (c *NeopixelCommand) Name() string {
	return c.fs.Name()
}

func (c *NeopixelCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *NeopixelCommand) Execute() error {
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

	neopixel := crickithat.NewNeopixel(dev)
	if err := neopixel.SetPin(uint8(c.pin)); err != nil {
		log.Fatalf("set pin: %w", err)
	}

	if err := neopixel.SetBufferLength(3); err != nil {
		log.Fatalf("set buffer length: %w", err)
	}

	buf := make([]byte, 3)
	// GRB
	buf[0] = byte(c.green)
	buf[1] = byte(c.red)
	buf[2] = byte(c.blue)

	if err := neopixel.SetBuffer(buf); err != nil {
		log.Fatalf("set buffer: %w", err)
	}

	if err := neopixel.Show(); err != nil {
		log.Fatalf("show: %w", err)
	}

	return nil
}
