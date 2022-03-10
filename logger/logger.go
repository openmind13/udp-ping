package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	LogLevelChan = make(chan string, 1)
)

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:04:05 02-01-2006",
		DisableColors:   false,
		FullTimestamp:   true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	go func() {
		for {
			level := <-LogLevelChan
			logrusLevel := logrus.GetLevel()
			switch level {
			case "trace":
				logrusLevel = logrus.TraceLevel
			case "debug":
				logrusLevel = logrus.DebugLevel
			case "info":
				logrusLevel = logrus.InfoLevel
			case "warn":
				logrusLevel = logrus.WarnLevel
			case "error":
				logrusLevel = logrus.ErrorLevel
			case "fatal":
				logrusLevel = logrus.FatalLevel
			default:
				logrus.Info("Unknown 'log_level': ", level, ". Using current log level")
			}
			if logrus.GetLevel() != logrusLevel {
				logrus.Info("Change log level to: ", logrusLevel)
				logrus.SetLevel(logrusLevel)
			}
		}
	}()
}
