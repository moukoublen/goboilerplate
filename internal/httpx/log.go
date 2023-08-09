package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type httpHeaderMarshalZerologObject struct {
	header http.Header
}

func (l *httpHeaderMarshalZerologObject) MarshalZerologObject(e *zerolog.Event) {
	for k := range l.header {
		e.Str(k, l.header.Get(k))
	}
}

type httpRequestMarshalZerologObject struct {
	Request *http.Request
}

func (rr *httpRequestMarshalZerologObject) MarshalZerologObject(e *zerolog.Event) {
	e.
		Object("header", &httpHeaderMarshalZerologObject{header: rr.Request.Header}).
		Int64("contentLength", rr.Request.ContentLength).
		Strs("transferEncoding", rr.Request.TransferEncoding).
		Str("host", rr.Request.Host).
		Str("remoteAddr", rr.Request.RemoteAddr).
		Str("uri", rr.Request.RequestURI).
		Str("proto", rr.Request.Proto).
		Str("method", rr.Request.Method).
		Str("url", requestURL(rr.Request))
}

type httpResponseMarshalZerologObject struct {
	Response *http.Response
}

func (rr *httpResponseMarshalZerologObject) MarshalZerologObject(e *zerolog.Event) {
	e.
		Object("header", &httpHeaderMarshalZerologObject{header: rr.Response.Header}).
		Str("status", rr.Response.Status).
		Int("statusCode", rr.Response.StatusCode).
		Str("proto", rr.Response.Proto).
		Int64("contentLength", rr.Response.ContentLength).
		Strs("transferEncoding", rr.Response.TransferEncoding)
}

type httpWrapResponseWriterMarshalZerologObject struct {
	ResponseWriter middleware.WrapResponseWriter
}

func (rr *httpWrapResponseWriterMarshalZerologObject) MarshalZerologObject(e *zerolog.Event) {
	e.
		Object("header", &httpHeaderMarshalZerologObject{header: rr.ResponseWriter.Header()}).
		Int("statusCode", rr.ResponseWriter.Status()).
		Int("bytesWritten", rr.ResponseWriter.BytesWritten())
}

type ChiZerolog struct {
	LogInLevel zerolog.Level
}

// NewLogEntry implements chi LogFormatter NewLogEntry function.
//
//nolint:ireturn
func (c *ChiZerolog) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := zerolog.Ctx(r.Context())

	if logger.GetLevel() > c.LogInLevel {
		return chiNopLogEntry{}
	}

	entry := &chiZerologEntry{
		request:  r,
		logEvent: logger.WithLevel(c.LogInLevel), //nolint:zerologlint // indented behavior.
		logger:   logger,
	}

	entry.logEvent.Str("request_id", middleware.GetReqID(r.Context()))
	entry.logEvent.EmbedObject(&httpRequestMarshalZerologObject{Request: r})

	return entry
}

type chiZerologEntry struct {
	request  *http.Request
	logEvent *zerolog.Event
	logger   *zerolog.Logger
}

func (e *chiZerologEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	e.logEvent.
		Int("status_code", status).
		Str("status", http.StatusText(status)).
		Int("bytes", bytes).
		Dur("elapsed", elapsed).
		Interface("extra", extra).
		Object("responseHeaders", &httpHeaderMarshalZerologObject{header: header}).
		Msg("http call")
}

func (e *chiZerologEntry) Panic(v any, stack []byte) {
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

func (e chiNopLogEntry) Write(_, _ int, _ http.Header, _ time.Duration, _ any) {}
func (e chiNopLogEntry) Panic(_ any, _ []byte)                                 {}

func requestURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto)
}

func dictFromWrapResponseWriter(w middleware.WrapResponseWriter, startTime time.Time, responseBodyBuffer bytes.Buffer) *zerolog.Event {
	dict := zerolog.Dict()
	dict.EmbedObject(&httpWrapResponseWriterMarshalZerologObject{ResponseWriter: w})
	dict.Dur("duration", time.Since(startTime))

	if logBody(w.Header()) {
		dict.RawJSON("payload", sanitizeJSONBytesToLog(responseBodyBuffer.Bytes()))
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
	dict.EmbedObject(&httpRequestMarshalZerologObject{Request: r})

	if logBody(r.Header) {
		dict.RawJSON("payload", sanitizeJSONBytesToLog(payloadBytes))
	}

	return dict, nil
}

func logBody(h http.Header) bool {
	contentType := h.Get("Content-Type")
	contentEncoding := h.Get("Content-Encoding")
	return strings.Contains(contentType, `application/json`) && len(contentEncoding) == 0
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
// Source: https://go.dev/src/net/http/httputil/dump.go
func drainBody(b io.ReadCloser) (io.ReadCloser, []byte, error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, nil, nil
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}

	if err := b.Close(); err != nil {
		return nil, nil, err
	}

	return io.NopCloser(&buf), buf.Bytes(), nil
}
