package io

type NoopWriter struct {
}

func NewNoopWriter() NoopWriter {
	return NoopWriter{}
}

func (w NoopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
