package crickithat

import (
	"encoding/binary"
)

const (
	SEESAW_NEOPIXEL_BASE       = 0x0e
	SEESAW_NEOPIXEL_PIN        = 0x01
	SEESAW_NEOPIXEL_SPEED      = 0x02
	SEESAW_NEOPIXEL_BUF_LENGTH = 0x03
	SEESAW_NEOPIXEL_BUF        = 0x04
	SEESAW_NEOPIXEL_SHOW       = 0x05
)

type Neopixel struct {
	d *Dev
}

func NewNeopixel(d *Dev) *Neopixel {
	return &Neopixel{d: d}
}

func (n *Neopixel) SetPin(pin uint8) error {
	data := make([]byte, 1)
	data[0] = pin

	err := n.d.writeRegister(SEESAW_NEOPIXEL_BASE, SEESAW_NEOPIXEL_PIN, data)
	return err
}

func (n *Neopixel) SetBufferLength(size uint16) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, size)

	err := n.d.writeRegister(SEESAW_NEOPIXEL_BASE, SEESAW_NEOPIXEL_BUF_LENGTH, data)
	return err
}

func (n *Neopixel) SetBuffer(buf []uint8) error {
	offset := []byte{0x00, 0x00}
	data := append(offset, buf...)

	err := n.d.writeRegister(SEESAW_NEOPIXEL_BASE, SEESAW_NEOPIXEL_BUF, data)
	return err
}

func (n *Neopixel) Show() error {
	err := n.d.writeRegister(SEESAW_NEOPIXEL_BASE, SEESAW_NEOPIXEL_SHOW, []byte{})
	return err
}
