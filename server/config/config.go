package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ListenUrl string `envconfig:"LISTEN_URL" default:":5000"`
}

func Load() *Config {
	var cfg Config
	err := envconfig.Process("SERVER", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &cfg
}
