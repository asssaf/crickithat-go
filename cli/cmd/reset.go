package cmd

import (
	"flag"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/crickithat-go/crickithat"
)

type ResetCommand struct {
	fs *flag.FlagSet
}

func NewResetCommand() *ResetCommand {
	c := &ResetCommand{
		fs: flag.NewFlagSet("reset", flag.ExitOnError),
	}

	return c
}

func (c *ResetCommand) Name() string {
	return c.fs.Name()
}

func (c *ResetCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *ResetCommand) Execute() error {
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

	return nil
}
