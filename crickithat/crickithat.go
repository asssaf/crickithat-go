package crickithat

import (
	"encoding/binary"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
)

const (
	SEESAW_STATUS_BASE    = 0x00
	SEESAW_STATUS_HW_ID   = 0x01
	SEESAW_STATUS_VERSION = 0x02
	SEESAW_STATUS_SWRST   = 0x7F

	SEESAW_TIMER_BASE = 0x08
	SEESAW_TIMER_PWM  = 0x01
	SEESAW_TIMER_FREQ = 0x02
)

type Opts struct {
	Addr uint16
}

var DefaultOpts = Opts{
	Addr: 0x49,
}

// Dev represents the device
type Dev struct {
	c    conn.Conn
	opts Opts
}

// NewI2C returns a new driver.
func NewI2C(b i2c.Bus, opts *Opts) (*Dev, error) {
	dev := &Dev{
		c:    &i2c.Dev{Bus: b, Addr: opts.Addr},
		opts: *opts,
	}

	return dev, nil
}

func (d *Dev) Init() error {
	return nil
}

// Halt all internal devices.
func (d *Dev) Halt() error {
	return nil
}

func (d *Dev) Reset() error {
	err := d.writeRegister(SEESAW_STATUS_BASE, SEESAW_STATUS_SWRST, []byte{0xff})
	return err
}

func (d *Dev) GetHardwareCode() (uint8, error) {
	data, err := d.readRegister(SEESAW_STATUS_BASE, SEESAW_STATUS_HW_ID, 1)
	if err != nil {
		return 0, err
	}

	return data[0], nil
}

func (d *Dev) GetVersion() (uint32, error) {
	data, err := d.readRegister(SEESAW_STATUS_BASE, SEESAW_STATUS_VERSION, 4)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data), nil
}

func (d *Dev) readRegister(addressHigh, addressLow uint8, length int) ([]byte, error) {
	write := []byte{addressHigh, addressLow}
	read := make([]byte, length)
	if err := d.c.Tx(write, read); err != nil {
		return nil, err
	}
	return read, nil
}

func (d *Dev) writeRegister(addressHigh, addressLow uint8, data []byte) error {
	write := []byte{addressHigh, addressLow}
	write = append(write, data...)
	if err := d.c.Tx(write, []byte{}); err != nil {
		return err
	}
	return nil
}
