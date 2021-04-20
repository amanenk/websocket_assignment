package logger

import (
	"go.uber.org/zap"
)

type RequestIdType string

var logger *zap.Logger

func Get() *zap.Logger {
	return logger
}

func Load() (err error) {
	logger, err = zap.NewProduction() // todo add log pattern here
	return err
}
