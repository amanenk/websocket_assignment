package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	//ListenUrl string `envconfig:"LISTEN_URL" default:":5000"`
	Port int `envconfig:"PORT" default:"5000"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
