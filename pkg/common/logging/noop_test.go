package logging

import "testing"

func TestNoopDontPanic(t *testing.T) {
	log := Noop()
	log.Debug("test %s", "ok")
	log.Info("test %s", "ok")
	log.Warn("test %s", "ok")
	log.Error("test %s", "ok")
}
