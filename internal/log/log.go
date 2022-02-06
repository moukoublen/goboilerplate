package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	if c.ConsoleLog {
		w = zerolog.ConsoleWriter{Out: os.Stdout}
	} else {
		w = os.Stdout
	}

	zerolog.SetGlobalLevel(c.LogLevel)
	log.Logger = newDefaultLogger(c.LogLevel, w)
	zerolog.DefaultContextLogger = &log.Logger
}

///////////////////////////////////
// Chi Log functionality

func LogRoutes(r *chi.Mux) {
	if log.Logger.GetLevel() > zerolog.DebugLevel {
		return
	}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Debug().Str("method", method).Str("route", route).Msg("log route")
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Error().Err(err).Msg("error during chi walk")
	}
}

type ChiZerolog struct {
	LogInLevel zerolog.Level
}

func (c *ChiZerolog) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := zerolog.Ctx(r.Context())

	if logger.GetLevel() > c.LogInLevel {
		return chiNopLogEntry{}
	}

	entry := &chiZerologEntry{
		request:  r,
		logEvent: logger.WithLevel(c.LogInLevel),
		logger:   logger,
	}

	entry.logEvent.Str("request_id", middleware.GetReqID(r.Context()))
	entry.logEvent.Str("method", r.Method)
	entry.logEvent.Str("url", requestURL(r))
	entry.logEvent.Str("uri", r.RequestURI)
	entry.logEvent.Str("from", r.RemoteAddr)
	entry.logEvent.Object("request_headers", headersLogObjectMarshaler{r.Header})

	return entry
}

type chiZerologEntry struct {
	request  *http.Request
	logEvent *zerolog.Event
	logger   *zerolog.Logger
}

func (e *chiZerologEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.logEvent.
		Int("status_code", status).
		Str("status", http.StatusText(status)).
		Int("bytes", bytes).
		Dur("elapsed", elapsed).
		Interface("extra", extra).
		Object("response_headers", headersLogObjectMarshaler{header}).
		Msg("http call")
}

func (e *chiZerologEntry) Panic(v interface{}, stack []byte) {
	s := strings.Split(string(stack), "\n")
	for i := range s {
		s[i] = strings.ReplaceAll(s[i], "\t", "  ")
	}

	e.logger.Error().
		Interface("panic", v).
		Strs("stack", s).
		Msg("http call paniced")
}

type chiNopLogEntry struct{}

func (e chiNopLogEntry) Write(_, _ int, _ http.Header, _ time.Duration, _ interface{}) {}
func (e chiNopLogEntry) Panic(_ interface{}, _ []byte)                                 {}

func requestURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto)
}

type headersLogObjectMarshaler struct {
	h http.Header
}

func (l headersLogObjectMarshaler) MarshalZerologObject(e *zerolog.Event) {
	for k := range l.h {
		e.Str(k, l.h.Get(k))
	}
}
