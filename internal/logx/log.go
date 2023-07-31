package logx

import (
	"io"
	"os"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func DefaultConfigValues() map[string]any {
	return map[string]any{
		"log.console_writer": false,
		"log.level":          0,
	}
}

type Config struct {
	ConsoleWriter bool
	LogLevel      zerolog.Level
}

func ParseConfig(cnf *koanf.Koanf) Config {
	return Config{
		ConsoleWriter: cnf.Bool("log.console_writer"),
		LogLevel:      zerolog.Level(cnf.Int("log.level")),
	}
}

func newDefaultLogger(l zerolog.Level, w io.Writer) zerolog.Logger {
	zerolog.SetGlobalLevel(l)
	zerolog.TimeFieldFormat = time.RFC3339

	//nolint: reassign
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// zerolog.ErrorMarshalFunc

	logger := zerolog.New(w).With().
		Stack().
		Timestamp().
		Caller().
		Logger()

	return logger
}

func SetupLog(c Config) {
	var w io.Writer
	if c.ConsoleWriter {
		w = zerolog.ConsoleWriter{Out: os.Stdout}
	} else {
		w = os.Stdout
	}

	zerolog.SetGlobalLevel(c.LogLevel)
	log.Logger = newDefaultLogger(c.LogLevel, w)
	zerolog.DefaultContextLogger = &log.Logger
}
