package main

import (
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"log"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
)

func main() {
	defer logger.FinalizeLogger()

	i2cHandle, err := i2c.NewI2C(0x40, 1)

	if err != nil {
		lg.Fatal(err)
	}
	defer i2cHandle.Close()
	logger.ChangePackageLogLevel("i2cHandle", logger.DebugLevel)

	htu21d := &HTU21D{i2cHandle: i2cHandle}

	_, err = htu21d.SoftRest()
	if err != nil {
		log.Fatal(err)
	}

	temp, err := htu21d.ReadTemp()
	if err != nil {
		log.Fatal(err)
	}

	humid, err := htu21d.ReadHumid()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Temp: %f, Humid: %f", temp, humid)
}
