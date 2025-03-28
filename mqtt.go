package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

var mq mqtt.Client

func setupMQTT() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("mqtt.broker"))
	opts.SetClientID(viper.GetString("mqtt.client_id"))
	opts.SetConnectTimeout(viper.GetDuration("mqtt.timeout"))
	opts.SetConnectRetryInterval(time.Second * 5)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Println("MQTT: connection lost:", err.Error())
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT: connected")
	})

	mq = mqtt.NewClient(opts)

	var (
		x        = mq.Connect()
		timedOut = x.WaitTimeout(time.Second)
	)

	if x.Error() != nil {
		return fmt.Errorf("MQTT: could not connect: %v", x.Error())
	} else if !timedOut {
		return fmt.Errorf("MQTT: connection timeout")
	}

	mq.Subscribe("kistan/kvitto", 0, handleMessage)

	return nil
}

func handleMessage(_ mqtt.Client, msg mqtt.Message) {
	if msg.Retained() {
		return
	}

	data := Kvitto{}
	data.Raw = msg.Payload()
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Println("MQTT: could not unmarshal payload:", err)
	}

	points := data.toPoints()
	for _, point := range points {
		if err := writeCli.WritePoint(context.Background(), point); err != nil {
			log.Println("Influx: could not write point:", err)
		}
	}
}
