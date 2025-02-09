package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileHook struct {
	Writer    *lumberjack.Logger
	Formatter logrus.Formatter
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLogger() {
	logrus.SetLevel(logrus.InfoLevel)

	consoleFormatter := &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	}

	fileFormatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	logFile := &lumberjack.Logger{
		Filename:   "fetcher.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(consoleFormatter)

	logrus.AddHook(&FileHook{
		Writer:    logFile,
		Formatter: fileFormatter,
	})
}
