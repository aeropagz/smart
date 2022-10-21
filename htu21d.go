package main

import (
	"encoding/binary"
	"github.com/d2r2/go-i2c"
	"time"
)

const (
	SOFT_RESET    = 0xFE
	READ_TEMP_NH  = 0xF3
	READ_HUMID_NH = 0xF5
	SLEEP_RESET   = 15 * time.Millisecond
	SLEEP_READ    = 50 * time.Millisecond
)

type HTU21D struct {
	i2cHandle *i2c.I2C
}

func (d *HTU21D) SoftRest() (int, error) {
	n, err := d.i2cHandle.WriteBytes([]byte{SOFT_RESET})
	if err != nil {
		return n, err
	}

	time.Sleep(SLEEP_RESET)
	return n, nil
}

func (d *HTU21D) ReadTemp() (float64, error) {
	_, err := d.triggerTemp()
	if err != nil {
		return 0.0, err
	}

	time.Sleep(SLEEP_READ)
	buf := make([]byte, 3)
	_, err = d.i2cHandle.ReadBytes(buf)
	if err != nil {
		return 0.0, err
	}

	out := binary.BigEndian.Uint16(buf)
	temp := -46.85 + (175.72 * (float64(out) / 65536.0))

	return temp, nil
}

func (d *HTU21D) ReadHumid() (float64, error) {
	_, err := d.triggerHumid()
	if err != nil {
		return 0, err
	}

	time.Sleep(SLEEP_READ)
	buf := make([]byte, 3)
	_, err = d.i2cHandle.ReadBytes(buf)
	if err != nil {
		return 0, err
	}

	out := binary.BigEndian.Uint16(buf)
	humid := -6.0 + (125.0 * (float64(out) / 65536.0))

	return humid, nil
}

func (d *HTU21D) triggerTemp() (int, error) {
	n, err := d.i2cHandle.WriteBytes([]byte{READ_TEMP_NH})
	if err != nil {
		return n, err
	}
	return n, nil
}

func (d *HTU21D) triggerHumid() (int, error) {
	n, err := d.i2cHandle.WriteBytes([]byte{READ_HUMID_NH})
	if err != nil {
		return n, err
	}
	return n, nil
}
