package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Url string `envconfig:"URL" default:"ws://localhost:5000/ws"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("CLIENT", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
