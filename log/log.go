package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// _defaultLevel is package default logging level.
	_defaultLevel = uint32(InfoLevel)
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func SetLevel(level Level) {
	_defaultLevel = uint32(level)
}

func Debugf(format string, args ...interface{}) {
	logf(DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	logf(InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logf(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logf(ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func logf(level Level, format string, args ...interface{}) {
	event := newEvent(level, format, args...)
	if uint32(event.Level) > _defaultLevel {
		return
	}

	switch level {
	case DebugLevel:
		logrus.WithTime(event.Time).Debugln(event.Message)
	case InfoLevel:
		logrus.WithTime(event.Time).Infoln(event.Message)
	case WarnLevel:
		logrus.WithTime(event.Time).Warnln(event.Message)
	case ErrorLevel:
		logrus.WithTime(event.Time).Errorln(event.Message)
	}
}
