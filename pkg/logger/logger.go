package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(logType, logFormat, logLevel string) *logrus.Logger {
	logger := logrus.New()

	if logType == "text" {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: logFormat,
			FullTimestamp:   true,
			ForceColors:     true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: logFormat,
		})
	}

	switch logLevel {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
