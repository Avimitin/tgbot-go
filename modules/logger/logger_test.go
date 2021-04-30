package logger

import "testing"

func TestLogger(t *testing.T) {
	log := NewZeroLogger("info")
	if log == nil {
		t.Fatal("logger is nil")
	}

	log.Info().Msg("this is a message")
}
