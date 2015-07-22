package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ErrBadLogLevel  = errors.New("bad log level")
	ErrBadLogFormat = errors.New("bad log format")
)

func setUpLogging() error {
	// level
	logrus.WithField("level", viper.GetString("log-level")).Debug("setting level")
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

	// format
	logrus.WithField("format", viper.GetString("log-format")).Debug("setting format")
	switch viper.GetString("log-format") {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		return ErrBadLogFormat
	}

	return nil
}
