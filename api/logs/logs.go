package logs

import (
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logrus.SetLevel(logrus.InfoLevel)

	fileFormatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	logFile := &lumberjack.Logger{
		Filename:   "api.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	logrus.SetOutput(logFile)
	logrus.SetFormatter(fileFormatter)
}
