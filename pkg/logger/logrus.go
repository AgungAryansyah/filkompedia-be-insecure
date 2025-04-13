package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    1,
		MaxAge:     7,
		MaxBackups: 5,
		Compress:   true,
	})

	return logger
}
