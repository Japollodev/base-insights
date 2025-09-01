package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// Parse log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	return logger
}
