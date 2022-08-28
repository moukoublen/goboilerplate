package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewHTTPInboundLoggerMiddleware(cnf config.Logging) func(http.Handler) http.Handler {
	switch cnf.InBoundHTTPLogLevel {
	case config.HTTPLogLevelNone:
		return func(h http.Handler) http.Handler { return h }
	case config.HTTPLogLevelBasic:
		return middleware.RequestLogger(&ChiZerolog{LogInLevel: cnf.LogInLevel})
	case config.HTTPLogLevelRequestResponse:
		return RequestResponseLogger(cnf)
	}

	return func(h http.Handler) http.Handler { return h }
}

func RequestResponseLogger(cnf config.Logging) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(r.Context()).With().Logger()

			requestDict, requestDictError := dictFromRequest(r)

			responseBodyBuffer := bytes.Buffer{}
			wrapResponseWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			wrapResponseWriter.Tee(&responseBodyBuffer)

			startTime := time.Now()

			next.ServeHTTP(wrapResponseWriter, r)

			responseDict := dictFromWrapResponseWriter(wrapResponseWriter, startTime, responseBodyBuffer)

			event := logger.WithLevel(cnf.LogInLevel)
			event.Str("url", r.URL.String())
			if requestDictError != nil {
				event.AnErr("requestDrainError", requestDictError)
			} else {
				event.Dict("request", requestDict)
			}
			event.Dict("response", responseDict)
			event.Msg("http transaction completed")
		})
	}
}

func dictFromWrapResponseWriter(w middleware.WrapResponseWriter, startTime time.Time, responseBodyBuffer bytes.Buffer) *zerolog.Event {
	dict := zerolog.Dict()
	dict.Interface("header", w.Header())
	dict.Int("status", w.Status())
	dict.Int("bytesWritten", w.BytesWritten())
	dict.Dur("duration", time.Since(startTime))

	if isContentTypeJSON(w.Header()) {
		dict.RawJSON("payload", sanitizeJSONBytesToLog(responseBodyBuffer.Bytes()))
	} else {
		dict.Str("payload", responseBodyBuffer.String())
	}

	return dict
}

func dictFromRequest(r *http.Request) (*zerolog.Event, error) {
	body, payloadBytes, err := drainBody(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = body

	dict := zerolog.Dict()
	dict.Interface("header", r.Header)
	dict.Int64("contentLength", r.ContentLength)
	dict.Strs("transferEncoding", r.TransferEncoding)
	dict.Str("host", r.Host)
	dict.Str("remoteAddr", r.RemoteAddr)
	dict.Str("RequestURI", r.RequestURI)
	dict.Str("proto", r.Proto)
	dict.Str("method", r.Method)

	if isContentTypeJSON(r.Header) {
		dict.RawJSON("payload", sanitizeJSONBytesToLog(payloadBytes))
	} else {
		dict.Str("payload", string(payloadBytes))
	}

	return dict, nil
}

func isContentTypeJSON(h http.Header) bool {
	contentType := h["Content-Type"]
	if len(contentType) > 0 {
		for _, h := range contentType {
			if strings.Contains(h, "application/json") {
				return true
			}
		}
	}

	return false
}

func sanitizeJSONBytesToLog(b []byte) []byte {
	return bytes.ReplaceAll(
		bytes.ReplaceAll(
			b,
			[]byte("\t"),
			[]byte{},
		),
		[]byte("\n"),
		[]byte{},
	)
}

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func drainBody(b io.ReadCloser) (r1 io.ReadCloser, r2 []byte, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, nil, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return io.NopCloser(&buf), buf.Bytes(), nil
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
