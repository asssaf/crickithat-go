package crickithat

import (
	"encoding/binary"
	"fmt"
)

const (
	PWM_SERVO4 = 0
	PWM_SERVO3 = 1
	PWM_SERVO2 = 2
	PWM_SERVO1 = 3

	MIN_PULSE = 3277
	MAX_PULSE = 6554

	FREQ_HZ = 50
)

var servoPwms = []uint8{PWM_SERVO1, PWM_SERVO2, PWM_SERVO3, PWM_SERVO4}

func (d *Dev) WriteServo(i int, value float64, minPulseMicros, maxPulseMicros int) error {
	if i < 0 || i > 3 {
		return fmt.Errorf("servo index should be in range 0-3: %d", i)
	}

	if value < 0.0 || value > 1.0 {
		return fmt.Errorf("value should be in range 0.0-1.0: %d", value)
	}

	// set pwm frequency
	if err := d.SetPwmFreq(i, FREQ_HZ); err != nil {
		return err
	}

	minPulse := float64(minPulseMicros) * FREQ_HZ * 65535 / 1_000_000
	if minPulse == 0 {
		minPulse = MIN_PULSE
	}

	maxPulse := float64(maxPulseMicros) * FREQ_HZ * 65535 / 1_000_000
	if maxPulse == 0 {
		maxPulse = MAX_PULSE
	}

	if minPulseMicros >= maxPulseMicros {
		return fmt.Errorf("the minimum pulse width must be less than the maximum pulse width")
	}

	scaledValue := uint16(scaleFloat64(value, 0.0, 1.0, minPulse, maxPulse))

	if err := d.SetWidth(i, scaledValue); err != nil {
		return err
	}

	return nil
}

func (d *Dev) WriteServoStop(i int) error {
	if i < 0 || i > 3 {
		return fmt.Errorf("servo index should be in range 0-3: %d", i)
	}

	if err := d.SetWidth(i, 0); err != nil {
		return err
	}

	return nil
}

func (d *Dev) SetPwmFreq(i int, frequency uint16) error {
	data := make([]byte, 3)
	data[0] = servoPwms[i]
	binary.BigEndian.PutUint16(data[1:], frequency)

	err := d.writeRegister(SEESAW_TIMER_BASE, SEESAW_TIMER_FREQ, data)
	return err
}

func (d *Dev) SetWidth(i int, value uint16) error {
	data := make([]byte, 3)
	data[0] = servoPwms[i]
	binary.BigEndian.PutUint16(data[1:], value)

	err := d.writeRegister(SEESAW_TIMER_BASE, SEESAW_TIMER_PWM, data)
	return err
}

func scaleFloat64(x, in_min, in_max, out_min, out_max float64) float64 {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}

func scale(x, in_min, in_max, out_min, out_max int) int {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}
