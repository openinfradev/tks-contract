package log

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

// Initialize initializes log.Logger with parameters.
func Initialize(serviceName string) {
	logger = log.New(os.Stdout, "["+serviceName+"] ", log.Ldate|log.Ltime)
}

func Logger() *log.Logger {
	return logger
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}
