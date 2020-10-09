package logging

import "github.com/sirupsen/logrus"

func Logrus() LogrusLogger {
	return LogrusLogger{}
}

type LogrusLogger struct {
}

func (l LogrusLogger) Debug(s string, i ...interface{}) {
	logrus.Debugf(s, i...)
}

func (l LogrusLogger) Info(s string, i ...interface{}) {
	logrus.Infof(s, i...)
}

func (l LogrusLogger) Warn(s string, i ...interface{}) {
	logrus.Warnf(s, i...)
}

func (l LogrusLogger) Error(s string, i ...interface{}) {
	logrus.Errorf(s, i...)
}
