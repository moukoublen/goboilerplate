package log

import (
	"io"
	"os"
	"time"

	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func newDefaultLogger(l zerolog.Level, w io.Writer) zerolog.Logger {
	zerolog.SetGlobalLevel(l)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// zerolog.ErrorMarshalFunc

	logger := zerolog.New(w).With().
		Stack().
		Timestamp().
		Caller().
		Logger()

	return logger
}

func SetupLog(c config.Logging) {
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
