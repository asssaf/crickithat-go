package cmd

import (
	"flag"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/crickithat-go/crickithat"
)

type ServoCommand struct {
	fs    *flag.FlagSet
	num   int
	value int
}

func NewServoCommand() *ServoCommand {
	c := &ServoCommand{
		fs: flag.NewFlagSet("servo", flag.ExitOnError),
	}

	c.fs.IntVar(&c.num, "num", 0, "Servo number (1-4)")
	c.fs.IntVar(&c.value, "value", -1, "Value to set (0-180)")

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
	if c.num < 1 || c.num > 4 {
		log.Fatalf("servo num must be in the range 1-4: %d", c.num)
	}

	if c.value < 0 || c.num > 180 {
		log.Fatalf("servo value must be in the range 0-180: %d", c.value)
	}

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

	if err := dev.Reset(); err != nil {
		log.Fatalf("device reset: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	// getting 0xff on the first read, just ignore it
	if _, err = dev.GetHardwareCode(); err != nil {
		log.Fatalf("read hardware code: %w", err)
	}

	expectedHwCode := uint8(0x55)
	if hwCode, err := dev.GetHardwareCode(); err != nil {
		log.Fatalf("read hardware code: %w", err)
	} else if hwCode != expectedHwCode {
		log.Printf("unexpected hardware code: 0x%x, expected 0x%x", hwCode, expectedHwCode)
	}

	if err := dev.WriteServo(c.num-1, c.value); err != nil {
		log.Fatalf("write servo: %w", err)
	}

	return nil
}
