package logs
/*
Package logs provides logging functionalities for the application.
It initializes and configures the logger to write logs in JSON format
to a file with rotation and compression settings.
*/

import (
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
InitLogger initializes the logger with specific settings.
It sets the log level to Info, formats the logs as JSON with a timestamp,
and configures the log output to a file with rotation and compression.
The log file is named "api.log", has a maximum size of 10 MB, keeps up to 3 backups,
and retains logs for 7 days.
*/

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
