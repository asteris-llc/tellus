package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ErrBadLogLevel = errors.New("bad log level")
)

func setUpLogging() error {
	logrus.WithField("level", viper.GetString("log-level")).Info("setting level")
	switch viper.GetString("log-level") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		return ErrBadLogLevel
	}

	return nil
}
