package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Entry {
	os.Mkdir("./logs", os.ModePerm)

	file, _ := os.OpenFile("./logs/logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}

	logger.Out = file

	contextLogger := logger.WithFields(logrus.Fields{
		"service": os.Getenv("DATABASE_SCHEMA"),
	})

	return contextLogger
}
