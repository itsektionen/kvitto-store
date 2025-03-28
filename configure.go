package main

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func configure() error {
	defaults()

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/it-sektionen")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("kvitto-store")

	if err := viper.ReadInConfig(); err != nil {
		var configFileAlreadyExistsError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileAlreadyExistsError) {
			// Write the default config to current directory if no file was found
			_ = viper.SafeWriteConfig()
		}

		return err
	}

	return nil
}

func defaults() {
	// MQTT settings
	viper.SetDefault("mqtt.broker", "localhost:1883")
	viper.SetDefault("mqtt.client_id", "kvitto-store")
	viper.SetDefault("mqtt.timeout", "1m")

	// Influx settings
	viper.SetDefault("influx.url", "http://localhost:8086")
	viper.SetDefault("influx.token", "suersecret")
	viper.SetDefault("influx.bucket", "kvitto")
	viper.SetDefault("influx.org", "smn")
}
