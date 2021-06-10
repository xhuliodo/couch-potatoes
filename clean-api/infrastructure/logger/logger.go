package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewAccessLogger() (logger *logrus.Logger) {
	accessLogFile := openLogFiles("./logs/access.log")
	logger = createLogger(accessLogFile)
	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

func NewErrorLogger() *ErrorLogger {
	errorLogFile := openLogFiles("./logs/error.log")
	logger := createLogger(errorLogFile)
	logger.Formatter = &logrus.JSONFormatter{}

	errorLogger := &ErrorLogger{logger}

	return errorLogger
}

func openLogFiles(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

func createLogger(file *os.File) *logrus.Logger {
	logger := logrus.New()
	logger.Out = file
	return logger
}
