package logger

import (
	"io"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerOutput struct {
	LogWriter     io.Writer
	LogType       string
	LogTimeFormat string
	LogLevel      string
}

func NewLogger(log LoggerOutput, hooks ...LoggerOutput) *logrus.Logger {
	logger := logrus.New()

	SetLogger(logger, log)
	for _, h := range hooks {
		SetHook(logger, h)
	}

	return logger
}

func SetLogger(logger *logrus.Logger, lo LoggerOutput) {
	if lo.LogType == "text" {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: lo.LogTimeFormat,
			FullTimestamp:   true,
			ForceColors:     true,
		})
	} else if lo.LogType == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: lo.LogTimeFormat,
		})
	}

	switch lo.LogLevel {
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
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetOutput(lo.LogWriter)
}

func SetHook(logger *logrus.Logger, lo LoggerOutput) {
	var formatter logrus.Formatter

	if lo.LogType == "text" {
		formatter = &logrus.TextFormatter{
			TimestampFormat: lo.LogTimeFormat,
			FullTimestamp:   true,
			ForceColors:     true,
		}
	} else if lo.LogType == "json" {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: lo.LogTimeFormat,
		}
	}

	writerMap := make(lfshook.WriterMap)

	switch lo.LogLevel {
	case "trace":
		writerMap[logrus.TraceLevel] = lo.LogWriter
		fallthrough
	case "debug":
		writerMap[logrus.DebugLevel] = lo.LogWriter
		fallthrough
	case "info":
		writerMap[logrus.InfoLevel] = lo.LogWriter
		fallthrough
	case "warn":
		writerMap[logrus.WarnLevel] = lo.LogWriter
		fallthrough
	case "error":
		writerMap[logrus.ErrorLevel] = lo.LogWriter
		fallthrough
	case "fatal":
		writerMap[logrus.FatalLevel] = lo.LogWriter
		fallthrough
	case "panic":
		writerMap[logrus.PanicLevel] = lo.LogWriter
	default:
		writerMap[logrus.InfoLevel] = lo.LogWriter
		writerMap[logrus.WarnLevel] = lo.LogWriter
		writerMap[logrus.ErrorLevel] = lo.LogWriter
		writerMap[logrus.FatalLevel] = lo.LogWriter
		writerMap[logrus.PanicLevel] = lo.LogWriter
	}

	logger.AddHook(lfshook.NewHook(
		writerMap,
		formatter,
	))
}

func NewLumberjackLogger(filePath string, fileMaxSize int) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename: filePath,
		MaxSize:  fileMaxSize, // MB
	}
}
