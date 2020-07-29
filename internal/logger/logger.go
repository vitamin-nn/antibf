package logger

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/config"
)

var (
	ErrEmptyLogLevel = errors.New("empty log level error")
)

func Init(logCfg config.Log) error {
	err := setLogLevel(logCfg.LogLevel)
	if err != nil {
		return err
	}
	return nil
}

func setLogLevel(logLevel string) error {
	if logLevel == "" {
		return ErrEmptyLogLevel
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}
