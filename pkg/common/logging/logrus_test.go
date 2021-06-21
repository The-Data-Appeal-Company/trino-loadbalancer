package logging

import "testing"

func TestLogrusDontPanic(t *testing.T) {
	log := LogrusLogger{}
	log.Debug("test %s", "ok")
	log.Info("test %s", "ok")
	log.Warn("test %s", "ok")
	log.Error("test %s", "ok")
}
