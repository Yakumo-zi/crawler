package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	plugin, c := NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := NewLogger(plugin)
	logger.Info("log init end")
}
