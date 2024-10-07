package httplog

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"time"
)

// RoundTripperFunc is an http.RoundTripper signature alias.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the RoundTripper interface.
func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func (il *HTTPLogger) LoggerRoundTripper(next http.RoundTripper) RoundTripperFunc {
	switch il.mode {
	case Drain:
		return il.loggerRoundTripperDrain(next)
	case Tee:
		return il.loggerRoundTripperTee(next)
	default:
		return il.loggerRoundTripperDrain(next)
	}
}

func (il *HTTPLogger) loggerRoundTripperDrain(next http.RoundTripper) RoundTripperFunc {
	return func(req *http.Request) (*http.Response, error) {
		attrs := make([]slog.Attr, 0, 4)

		attrs = append(attrs, il.attrConverter.HTTPRequest(req))

		startTime := time.Now()

		// next transport
		res, err := next.RoundTrip(req)
		attrs = append(attrs, slog.Duration("duration", time.Since(startTime)))
		if err != nil {
			attrs = append(attrs, attrError("sendError", err))
		}

		attrs = append(attrs, il.attrConverter.HTTPResponse(res))

		il.logOutbound(req.Context(), attrs)

		return res, err
	}
}

func (il *HTTPLogger) loggerRoundTripperTee(next http.RoundTripper) RoundTripperFunc {
	return func(req *http.Request) (*http.Response, error) {
		reqAttrs := il.attrConverter.AttrsHTTPRequestExcludeBody(req)

		req.Body = NewTeeReadCloserPooled(req.Body, il.pool, func(readErr, closeErr error, buf *bytes.Buffer) {
			if readErr != nil {
				reqAttrs = append(reqAttrs, attrError("readError", readErr))
			}
			if closeErr != nil {
				reqAttrs = append(reqAttrs, attrError("closeError", closeErr))
			}
			reqAttrs = append(reqAttrs, attrBody(buf.Bytes()))
		})

		startTime := time.Now()

		// next transport
		res, err := next.RoundTrip(req)

		attrs := make([]slog.Attr, 0, 4)
		attrs = append(attrs, slog.Duration("duration", time.Since(startTime)))
		if err != nil {
			attrs = append(attrs, attrError("sendError", err))
		}

		attrs = append(attrs, il.attrConverter.GroupAttrsAsHTTPRequest(reqAttrs)) // request

		resAttrs := il.attrConverter.AttrsHTTPResponseExcludeBody(res)
		res.Body = NewTeeReadCloserPooled(res.Body, il.pool, func(readErr, closeErr error, buf *bytes.Buffer) {
			if readErr != nil {
				resAttrs = append(resAttrs, attrError("readError", readErr))
			}
			if closeErr != nil {
				resAttrs = append(resAttrs, attrError("closeError", closeErr))
			}
			resAttrs = append(resAttrs, attrBody(buf.Bytes()))

			attrs = append(attrs, il.attrConverter.GroupAttrsAsHTTPResponse(resAttrs)) // response

			il.logOutbound(req.Context(), attrs)
		})

		return res, err
	}
}

func (il *HTTPLogger) logOutbound(ctx context.Context, attrs []slog.Attr) {
	il.logger.LogAttrs(ctx, il.logInLevel.Level(), "http outbound", attrs...)
}
