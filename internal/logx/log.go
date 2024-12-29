package logx

import (
	"context"
	"log/slog"
	"os"

	"github.com/knadh/koanf/v2"
)

type LogType string

const (
	LogTypeJSON LogType = "json"
	LogTypeText LogType = "text"
)

type Config struct {
	LogType LogType
	Level   slog.Level
}

func ParseConfig(cnf *koanf.Koanf) Config {
	level := slog.Level(0)
	levelStr := cnf.Bytes("log.level")
	err := level.UnmarshalText(levelStr)
	_ = err

	return Config{
		Level:   level,
		LogType: LogType(cnf.String("log.type")),
	}
}

func DefaultConfigValues() map[string]any {
	return map[string]any{
		"log.type":  "text",
		"log.level": "INFO",
	}
}

//nolint:gochecknoglobals
var logLevel = &slog.LevelVar{}

func SetLogLevel(l slog.Level) {
	slog.SetLogLoggerLevel(l)
	logLevel.Set(l)
}

func InitSLog(c Config, attrs ...slog.Attr) *slog.Logger {
	SetLogLevel(c.Level)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	var slogHandler slog.Handler
	switch c.LogType { //nolint:exhaustive
	case LogTypeJSON:
		slogHandler = slog.NewJSONHandler(os.Stdout, opts)
	default: // fallback to text
		slogHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(slogHandler.WithAttrs(attrs))

	slog.SetDefault(logger)

	logger.Debug("logger initialized", slog.String("level", c.Level.String()))

	return logger
}

type ctxSLogKey struct{}

func GetFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxSLogKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func SetInContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxSLogKey{}, logger)
}

type NOOPLogHandler struct{}

func (NOOPLogHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (NOOPLogHandler) Handle(context.Context, slog.Record) error { return nil }
func (l NOOPLogHandler) WithAttrs(_ []slog.Attr) slog.Handler    { return l }
func (l NOOPLogHandler) WithGroup(_ string) slog.Handler         { return l }

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	return slog.String("error", err.Error())
}

func AnError(key string, err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	return slog.String(key, err.Error())
}
