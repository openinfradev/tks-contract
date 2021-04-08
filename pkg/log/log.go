package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

// Init initializes logrus.logger and set
func init() {
	logger = logrus.New()
	logger.Out = os.Stdout

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Error(v ...interface{}) {
	logger.Debug(v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
