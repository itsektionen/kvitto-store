package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting kvitto-store...")

	log.Println("Loading config...")
	if err := configure(); err != nil {
		log.Fatal(err)
	}

	log.Println("Setting up MQTT...")
	if err := setupMQTT(); err != nil {
		log.Fatal(err)
	}

	log.Println("Setting up InfluxDB...")
	setupInflux()

	// Setup and wait for exit signal
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wait for exit signal
	<-sigChannel
	mq.Disconnect(uint(time.Second.Milliseconds()))
}
