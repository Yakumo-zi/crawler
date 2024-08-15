package log

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
	DPanicLevel = zapcore.DPanicLevel
)

type Plugin zapcore.Core

func NewStdoutPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stdout)), enabler)
}

func NewStderrPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stderr)), enabler)
}

func NewFilePlugin(
	filePath string, enabler zapcore.LevelEnabler) (Plugin, io.Closer) {
	var writer = DefaultLumberjackLogger()
	writer.Filename = filePath
	return NewPlugin(zapcore.AddSync(writer), enabler), writer
}

func NewPlugin(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) Plugin {
	return zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writer, enabler)
}

func DefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:   200,
		LocalTime: true,
		Compress:  true,
	}
}

func NewLogger(plugin zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(plugin, append(DefaultOption(), options...)...)
}

func DefaultOption() []zap.Option {
	var stackTraceLevel zap.LevelEnablerFunc = func(l zapcore.Level) bool {
		return l >= zapcore.DPanicLevel
	}
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(stackTraceLevel),
	}
}
