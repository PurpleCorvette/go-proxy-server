// Package logger provides a configured logrus logger.
package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	loglevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(loglevel)
	}

	return logger
}
