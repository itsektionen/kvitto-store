package main

import (
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/viper"
)

var (
	cli      influxdb2.Client
	writeCli api.WriteAPIBlocking
)

func setupInflux() {
	cli = influxdb2.NewClient(viper.GetString("influx.url"), viper.GetString("influx.token"))
	writeCli = cli.WriteAPIBlocking(viper.GetString("influx.organization"), viper.GetString("influx.bucket"))
}
