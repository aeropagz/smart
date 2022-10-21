package main

import (
	"context"
	"log"
	"time"

	htu21d2 "github.com/aeropagz/smart/htu21d"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/influxdata/influxdb-client-go/v2"
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

	htu21d := &htu21d2.HTU21D{I2cHandle: i2cHandle, SensorName: "Schlafzimmer"}

	_, err = htu21d.SoftRest()
	if err != nil {
		log.Fatal(err)
	}

	result, err := htu21d.GetResult()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s	Temp: %f, Humid: %f", result.SensorName, result.Temp, result.Humid)

	client := influxdb2.NewClient("http://localhost:8086", "super-secret-token")
	writeAPI := client.WriteAPIBlocking("my-org", "smart")
	p := influxdb2.NewPointWithMeasurement("sensor").
		AddTag("sensor-name", result.SensorName).
		AddField("temp", result.Temp).
		AddField("humid", result.Humid).
		SetTime(time.Now())
	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}
