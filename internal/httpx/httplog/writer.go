package httplog

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

type netHTTPResponse interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	FlushError() error                           // *http.response
	ReadFrom(src io.Reader) (n int64, err error) // *http.response
	SetReadDeadline(deadline time.Time) error    // *http.response
	SetWriteDeadline(deadline time.Time)         // *http.response
	EnableFullDuplex() error                     // *http.response
	WriteString(data string) (n int, err error)  // *http.response
}

type ResponseWriterWrapper interface {
	http.ResponseWriter
	Header() http.Header
	Buffer() *bytes.Buffer
	Status() int
	BytesWritten() int
}

func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper { //nolint:ireturn,cyclop
	asFlusher, isFlusher := w.(http.Flusher)
	asPusher, isPusher := w.(http.Pusher)
	asReaderFrom, isReaderFrom := w.(io.ReaderFrom)
	asHijacker, isHijacker := w.(http.Hijacker)
	asNetHTTPResponse, isNetHTTPResponse := w.(netHTTPResponse)

	var wrapperResponseWriter *responseWriterWrapper
	var wrapperFlusher *flusher
	var wrapperPusher *pusher
	var wrapperReaderFrom *readerFrom
	var wrapperHijacker *hijacker
	var wrapperNetHTTPResponse *netHTTPResponseWrapper

	wrapperResponseWriter = &responseWriterWrapper{
		wrapped: w,
		tee:     &bytes.Buffer{},
	}

	if isFlusher {
		wrapperFlusher = &flusher{w: wrapperResponseWriter, f: asFlusher}
	}

	if isPusher {
		wrapperPusher = &pusher{p: asPusher}
	}

	if isReaderFrom {
		wrapperReaderFrom = &readerFrom{w: wrapperResponseWriter, r: asReaderFrom}
	}

	if isHijacker {
		wrapperHijacker = &hijacker{h: asHijacker}
	}

	if isNetHTTPResponse {
		wrapperNetHTTPResponse = &netHTTPResponseWrapper{
			wrapped:               asNetHTTPResponse,
			responseWriterWrapper: wrapperResponseWriter,
			readerFrom:            wrapperReaderFrom,
			flusher:               wrapperFlusher,
			hijacker:              wrapperHijacker,
		}
	}

	switch {
	// [0000] simple http.ResponseWriter interface
	case !isFlusher || !isPusher || !isReaderFrom || !isHijacker || !isNetHTTPResponse:
		return wrapperResponseWriter

	// *http.response with or without pusher
	case !isPusher && isNetHTTPResponse:
		return wrapperNetHTTPResponse
	case isPusher && isNetHTTPResponse:
		return netHTTPResponseAndPusherWrapper{netHTTPResponseWrapper: wrapperNetHTTPResponse, pusher: wrapperPusher}

	// [0001] http.ResponseWriter + http.Hijacker
	case !isFlusher && !isPusher && !isReaderFrom && isHijacker:
		return responseWriterAndHijacker{responseWriterWrapper: wrapperResponseWriter, hijacker: wrapperHijacker}

	// [0010] http.ResponseWriter + io.ReaderFrom
	case !isFlusher && !isPusher && isReaderFrom && !isHijacker:
		return responseWriterAndReaderFrom{responseWriterWrapper: wrapperResponseWriter, readerFrom: wrapperReaderFrom}

	// TODO: [0011] !isFlusher && !isPusher && isReaderFrom && isHijacker

	// [0100] http.ResponseWriter + http.isPusher
	case !isFlusher && isPusher && !isReaderFrom && !isHijacker:
		return responseWriterAndPusher{responseWriterWrapper: wrapperResponseWriter, pusher: wrapperPusher}

	// TODO: [0101] !isFlusher && isPusher && !isReaderFrom && isHijacker
	// TODO: [0110] !isFlusher && isPusher && isReaderFrom && !isHijacker
	// TODO: [0111] !isFlusher && isPusher && isReaderFrom && isHijacker

	// [1000] http.ResponseWriter + http.Flusher
	case isFlusher && !isPusher && !isReaderFrom && !isHijacker:
		return responseWriterAndFlusher{responseWriterWrapper: wrapperResponseWriter, flusher: wrapperFlusher}

	// TODO: [1001] isFlusher && !isPusher && !isReaderFrom && isHijacker
	// TODO: [1010] isFlusher && !isPusher && isReaderFrom && !isHijacker
	// TODO: [1011] isFlusher && !isPusher && isReaderFrom && isHijacker

	// [1100] http.ResponseWriter + http.Flusher + http.Pusher
	case isFlusher && isPusher && !isReaderFrom && !isHijacker:
		return responseWriterFlusherPusher{responseWriterWrapper: wrapperResponseWriter, flusher: wrapperFlusher, pusher: wrapperPusher}

	// TODO: [1101] isFlusher && isPusher && !isReaderFrom && isHijacker
	// TODO: [1110] isFlusher && isPusher && isReaderFrom && !isHijacker
	// TODO: [1111] isFlusher && isPusher && isReaderFrom && isHijacker

	default:
		return wrapperResponseWriter
	}
}

// ### http.ResponseWriter + http.Flusher + http.isPusher
var (
	_ http.ResponseWriter = (*responseWriterFlusherPusher)(nil)
	_ http.Flusher        = (*responseWriterFlusherPusher)(nil)
	_ http.Pusher         = (*responseWriterFlusherPusher)(nil)
)

type responseWriterFlusherPusher struct {
	*responseWriterWrapper
	*flusher
	*pusher
}

// ### http.ResponseWriter + http.Hijacker
var (
	_ http.ResponseWriter = (*responseWriterAndHijacker)(nil)
	_ http.Hijacker       = (*responseWriterAndHijacker)(nil)
)

type responseWriterAndHijacker struct {
	*responseWriterWrapper
	*hijacker
}

// ### http.ResponseWriter + io.ReaderFrom
var (
	_ http.ResponseWriter = (*responseWriterAndReaderFrom)(nil)
	_ io.ReaderFrom       = (*responseWriterAndReaderFrom)(nil)
)

type responseWriterAndReaderFrom struct {
	*responseWriterWrapper
	*readerFrom
}

// ### http.ResponseWriter + http.Pusher
var (
	_ http.ResponseWriter = (*responseWriterAndPusher)(nil)
	_ http.Pusher         = (*responseWriterAndPusher)(nil)
)

type responseWriterAndPusher struct {
	*responseWriterWrapper
	*pusher
}

// ### http.ResponseWriter + http.isFlusher
var (
	_ http.ResponseWriter = (*responseWriterAndFlusher)(nil)
	_ http.Flusher        = (*responseWriterAndFlusher)(nil)
)

type responseWriterAndFlusher struct {
	*responseWriterWrapper
	*flusher
}

// ### http.ResponseWriter
var _ http.ResponseWriter = (*responseWriterWrapper)(nil)

type responseWriterWrapper struct {
	wrapped     http.ResponseWriter
	tee         *bytes.Buffer
	statusCode  int
	bytes       int
	wroteHeader bool
}

func (w *responseWriterWrapper) Unwrap() http.ResponseWriter {
	return w.wrapped
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.statusCode = statusCode
		w.wroteHeader = true
		w.wrapped.WriteHeader(statusCode)
	}
}

func (w *responseWriterWrapper) Write(buf []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	if w.tee == nil {
		return w.wrapped.Write(buf)
	}

	n, err = w.wrapped.Write(buf)
	if w.tee != nil {
		_, teeErr := w.tee.Write(buf[:n])
		err = errors.Join(err, teeErr)
	}

	w.bytes += n

	return n, err
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.wrapped.Header()
}

func (w *responseWriterWrapper) Status() int {
	return w.statusCode
}

func (w *responseWriterWrapper) BytesWritten() int {
	return w.bytes
}

func (w *responseWriterWrapper) Buffer() *bytes.Buffer {
	return w.tee
}

// ### http.Flusher
var _ http.Flusher = (*flusher)(nil)

type flusher struct {
	w *responseWriterWrapper
	f http.Flusher
}

func (f *flusher) Flush() {
	f.w.wroteHeader = true
	f.f.Flush()
}

// ### http.Hijacker
var _ http.Hijacker = (*hijacker)(nil)

type hijacker struct {
	h http.Hijacker
}

func (h *hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.h.Hijack()
}

// ### http.Pusher
var _ http.Pusher = (*pusher)(nil)

type pusher struct {
	p http.Pusher
}

func (p *pusher) Push(target string, opts *http.PushOptions) error {
	return p.p.Push(target, opts)
}

// ### io.ReaderFrom
var _ io.ReaderFrom = (*readerFrom)(nil)

type readerFrom struct {
	w *responseWriterWrapper
	r io.ReaderFrom //nolint:unused
}

func (w *readerFrom) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.Copy(w.w, r) // tunnel to responseWriterWrapper.Write
	w.w.bytes += int(n)
	return n, err
}

// ### *net.response
var _ netHTTPResponse = (*netHTTPResponseWrapper)(nil)

type netHTTPResponseWrapper struct {
	wrapped netHTTPResponse
	*responseWriterWrapper
	*readerFrom
	*flusher
	*hijacker
}

func (w *netHTTPResponseWrapper) FlushError() error {
	return w.wrapped.FlushError()
}

func (w *netHTTPResponseWrapper) SetReadDeadline(deadline time.Time) error {
	return w.wrapped.SetReadDeadline(deadline)
}

func (w *netHTTPResponseWrapper) SetWriteDeadline(deadline time.Time) {
	w.wrapped.SetWriteDeadline(deadline)
}

func (w *netHTTPResponseWrapper) EnableFullDuplex() error {
	return w.wrapped.EnableFullDuplex()
}

func (w *netHTTPResponseWrapper) WriteString(data string) (n int, err error) {
	return w.Write([]byte(data)) // tunnel to write
}

// ### *net.response + http.Pusher
var (
	_ netHTTPResponse = (*netHTTPResponseAndPusherWrapper)(nil)
	_ http.Pusher     = (*netHTTPResponseAndPusherWrapper)(nil)
)

type netHTTPResponseAndPusherWrapper struct {
	*netHTTPResponseWrapper
	*pusher
}
