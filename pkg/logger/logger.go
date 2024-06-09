// Package logger provides a configured logrus logger.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger defines an interface for logging.
type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// LogrusLogger is an implementation of Logger using logrus.
type LogrusLogger struct {
	*logrus.Logger
}

// NewLogger initializes a new logrus logger with the specified log level.
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
