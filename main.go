package main

import (
	"encoding/binary"
	"time"

	i2c "github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
)

const (
	SOFT_RESET    = 0xFE
	READ_TEMP_NH  = 0xF3
	READ_HUMID_NH = 0xF5
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
)

func main() {
	defer logger.FinalizeLogger()

	i2c, err := i2c.NewI2C(0x40, 1)
	if err != nil {
		lg.Fatal(err)
	}
	defer i2c.Close()

	logger.ChangePackageLogLevel("i2c", logger.DebugLevel)

	n, err := i2c.WriteBytes([]byte{SOFT_RESET})
	if err != nil {
		lg.Fatal(err)
	}
	if n == 0 {
		lg.Info("Soft reset failed")
	}

	time.Sleep(50000 * time.Microsecond)
	n, err = i2c.WriteBytes([]byte{READ_TEMP_NH})
	if err != nil {
		lg.Fatal(err)
	}
	if n == 0 {
		lg.Info("Starting Measurement failed")
	}

	i := 0
	var result uint16 = 0
	for i < 5 {
		time.Sleep(50000 * time.Microsecond)
		bytesRead := make([]byte, 3)
		n, err = i2c.ReadBytes(bytesRead)
		if err != nil {
			lg.Fatal(err)
		}
		if n == 0 {
			i++
			lg.Info("nothing read, try: ", i)
		} else {
			result = binary.BigEndian.Uint16(bytesRead)
			break
		}
	}
	temp := -46.85 + (175.72 * (float64(result) / 65536.0))
	lg.Infof("temp: %d", temp)
}
