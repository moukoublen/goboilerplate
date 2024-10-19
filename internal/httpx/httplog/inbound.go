package httplog

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"time"
)

func (il *HTTPLogger) Handler(next http.Handler) http.Handler {
	switch il.mode {
	case Drain:
		return il.handlerDrain(next)
	case Tee:
		return il.handlerTee(next)
	default:
		return il.handlerDrain(next)
	}
}

func (il *HTTPLogger) handlerDrain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAttr := il.attrConverter.HTTPRequest(r)

		wrapResponseWriter := NewResponseWriterWrapper(w)

		startTime := time.Now()

		// serve
		next.ServeHTTP(wrapResponseWriter, r)

		// prepare attrs
		attrs := make([]slog.Attr, 0, 10)

		attrs = append(attrs, slog.Duration("duration", time.Since(startTime)))

		responseAttr := il.attrConverter.HTTPResponseWriter(wrapResponseWriter.Header(), wrapResponseWriter.Status(), wrapResponseWriter.Buffer().Bytes())

		attrs = append(attrs, requestAttr)
		attrs = append(attrs, responseAttr)

		il.logInbound(r.Context(), attrs)
	})
}

func (il *HTTPLogger) handlerTee(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqAttrs := il.attrConverter.AttrsHTTPRequestExcludeBody(r)

		r.Body = NewTeeReadCloserPooled(r.Body, il.pool, func(readErr, closeErr error, buf *bytes.Buffer) {
			if readErr != nil {
				reqAttrs = append(reqAttrs, attrError("readError", readErr))
			}
			if closeErr != nil {
				reqAttrs = append(reqAttrs, attrError("closeError", closeErr))
			}
			reqAttrs = append(reqAttrs, attrBody(buf.Bytes()))
		})

		wrapResponseWriter := NewResponseWriterWrapper(w)

		startTime := time.Now()

		// serve
		next.ServeHTTP(wrapResponseWriter, r)

		// prepare attrs
		attrs := make([]slog.Attr, 0, 3)
		attrs = append(attrs, slog.Duration("duration", time.Since(startTime)))
		attrs = append(attrs, il.attrConverter.GroupAttrsAsHTTPRequest(reqAttrs))
		attrs = append(
			attrs,
			il.attrConverter.HTTPResponseWriter(wrapResponseWriter.Header(), wrapResponseWriter.Status(), wrapResponseWriter.Buffer().Bytes()),
		)

		il.logInbound(r.Context(), attrs)
	})
}

func (il *HTTPLogger) logInbound(ctx context.Context, attrs []slog.Attr) {
	il.logger.LogAttrs(ctx, il.logInLevel.Level(), "http inbound", attrs...)
}
