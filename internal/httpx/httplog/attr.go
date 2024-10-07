package httplog

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPSLogAttrsConverter struct {
	logPolicy LogPolicy
}

func (a HTTPSLogAttrsConverter) Headers(key string, h http.Header) slog.Attr {
	s := make([]slog.Attr, 0, len(h))

	for key, value := range h {
		if a.logPolicy.ShouldOmitHeader(key, value) {
			continue
		}
		if a.logPolicy.ShouldMaskHeader(key, value) {
			s = append(s, slog.String(key, "***"))
			continue
		}
		s = append(s, slog.String(key, h.Get(key)))
	}

	return slog.Attr{Key: key, Value: slog.GroupValue(s...)}
}

func (a HTTPSLogAttrsConverter) URL(u *url.URL) slog.Attr {
	s := make([]slog.Attr, 0, 6)

	s = append(s, slog.String("scheme", u.Scheme))
	s = append(s, slog.String("opaque", u.Opaque))
	s = append(s, slog.String("host", u.Host))
	s = append(s, slog.String("path", u.Path))
	s = append(s, slog.String("fragment", u.Fragment))
	s = append(s, slog.String("full", fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)))

	return slog.Attr{Key: "url", Value: slog.GroupValue(s...)}
}

func (a HTTPSLogAttrsConverter) HTTPRequest(r *http.Request) slog.Attr {
	s := a.AttrsHTTPRequestExcludeBody(r)

	s = append(s, a.AttrsHTTPRequestBodyDrain(r)...)

	return a.GroupAttrsAsHTTPRequest(s)
}

func (a HTTPSLogAttrsConverter) AttrsHTTPRequestExcludeBody(r *http.Request) []slog.Attr {
	s := make([]slog.Attr, 0, 16)

	s = append(s, a.URL(r.URL))
	s = append(s, a.Headers("headers", r.Header))
	s = append(s, slog.String("method", r.Method))
	s = append(s, slog.String("proto", r.Proto))
	s = append(s, slog.Int64("contentLength", r.ContentLength))

	if len(r.TransferEncoding) > 0 {
		s = append(s, slog.Any("transferEncoding", r.TransferEncoding))
	}

	s = append(s, slog.String("host", r.Host))

	// Form - for inbound: This field is only available after ParseForm is called.
	if len(r.Form) > 0 {
		s = append(s, attrURLValues("form", r.Form))
	}

	// PostForm - for inbound: This field is only available after ParseForm is called.
	if len(r.PostForm) > 0 {
		s = append(s, attrURLValues("postForm", r.PostForm))
	}

	// MultipartForm - for inbound: This field is only available after ParseMultipartForm
	// if r.MultipartForm != nil {}

	if len(r.Trailer) > 0 {
		s = append(s, a.Headers("trailerHeaders", r.Trailer))
	}

	s = append(s, slog.String("remoteAddr", r.RemoteAddr))
	s = append(s, slog.String("requestUri", r.RequestURI))

	// TLS *tls.ConnectionState
	if r.TLS != nil {
		s = append(s, attrTLS("tls", r.TLS))
	}

	// Pattern string -> go 1.23
	// s = append(s, slog.String("requestUri", r.Pattern))

	return s
}

func (a HTTPSLogAttrsConverter) AttrsHTTPRequestBodyDrain(r *http.Request) []slog.Attr {
	s := make([]slog.Attr, 0, 2)

	// body dump
	switch {
	case r.Body == nil || r.Body == http.NoBody:
		s = append(s, slog.String("bodyLogNote", "no body"))
	case !a.logPolicy.ShouldLogRequestBody(r):
		s = append(s, slog.String("bodyLogNote", "body is not logable"))
	default:
		payloadBytes, err := drainRequestBody(r)
		if err != nil {
			s = append(s, slog.String("bodyLogNote", "log body error - %s"+err.Error()))
		}
		s = append(s, attrBody(payloadBytes))
	}

	return s
}

func (a HTTPSLogAttrsConverter) GroupAttrsAsHTTPRequest(s []slog.Attr) slog.Attr {
	return slog.Attr{Key: "request", Value: slog.GroupValue(s...)}
}

func (a HTTPSLogAttrsConverter) HTTPResponse(r *http.Response) slog.Attr {
	s := a.AttrsHTTPResponseExcludeBody(r)

	s = append(s, a.AttrsHTTPResponseDrainBody(r)...)

	return a.GroupAttrsAsHTTPResponse(s)
}

func (a HTTPSLogAttrsConverter) AttrsHTTPResponseExcludeBody(r *http.Response) []slog.Attr {
	s := make([]slog.Attr, 0, 16)

	s = append(s, attrStatusCode(r.StatusCode))
	s = append(s, slog.String("proto", r.Proto))
	s = append(s, a.Headers("headers", r.Header))
	s = append(s, slog.Int64("contentLength", r.ContentLength))
	if len(r.TransferEncoding) > 0 {
		s = append(s, slog.Any("transferEncoding", r.TransferEncoding))
	}
	s = append(s, slog.Bool("uncompressed", r.Uncompressed))
	if len(r.Trailer) > 0 {
		s = append(s, a.Headers("trailerHeaders", r.Trailer))
	}

	if r.TLS != nil {
		s = append(s, attrTLS("tls", r.TLS))
	}

	return s
}

func (a HTTPSLogAttrsConverter) AttrsHTTPResponseDrainBody(r *http.Response) []slog.Attr {
	s := make([]slog.Attr, 0, 2)

	switch {
	case r.Body == nil || r.Body == http.NoBody:
		s = append(s, slog.String("bodyLogNote", "no body"))
	case !a.logPolicy.ShouldLogResponseBody(r):
		s = append(s, slog.String("bodyLogNote", "body is not logable"))
	default:
		payloadBytes, err := drainResponseBody(r)
		if err != nil {
			s = append(s, slog.String("bodyLogNote", "log body error - "+err.Error()))
		}
		s = append(s, attrBody(payloadBytes))
	}

	return s
}

func (a HTTPSLogAttrsConverter) GroupAttrsAsHTTPResponse(s []slog.Attr) slog.Attr {
	return slog.Attr{Key: "response", Value: slog.GroupValue(s...)}
}

type BytesBuffer interface {
	Bytes() []byte
}

func (a HTTPSLogAttrsConverter) HTTPResponseWriter(headers http.Header, statusCode int, body []byte) slog.Attr {
	s := make([]slog.Attr, 0, 3)

	s = append(s, a.Headers("headers", headers))
	s = append(s, attrStatusCode(statusCode))

	if a.logPolicy.ShouldLogResponseWriterBody(headers, statusCode, body) {
		s = append(s, attrBody(body))
	}

	return slog.Attr{Key: "response", Value: slog.GroupValue(s...)}
}

func attrError(key string, err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	return slog.String(key, err.Error())
}

func attrStatusCode(statusCode int) slog.Attr {
	return slog.Group("status", slog.Int("code", statusCode), slog.String("name", http.StatusText(statusCode)))
}

func attrURLValues(key string, u url.Values) slog.Attr {
	if len(u) == 0 {
		return slog.Attr{}
	}

	attrs := make([]slog.Attr, 0, len(u))
	for k := range u {
		attrs = append(attrs, slog.String(k, u.Get(k)))
	}

	return slog.Attr{
		Key:   key,
		Value: slog.GroupValue(attrs...),
	}
}

func attrTLS(key string, tls *tls.ConnectionState) slog.Attr {
	return slog.Group(
		key,
		slog.Uint64("version", uint64(tls.Version)),
		slog.String("negotiatedProtocol", tls.NegotiatedProtocol),
	)
}

func attrBody(body []byte) slog.Attr {
	return slog.Group(
		"body",
		slog.Int("size", len(body)),
		slog.String("value", sanitizeJSONBytesToLog(body)),
	)
}

func sanitizeJSONBytesToLog(b []byte) string {
	s := strconv.QuoteToGraphic(string(b))

	if (len(s) > 2) && (s[0] == '"') && (s[len(s)-1] == '"') {
		return s[1 : len(s)-1]
	}

	return strconv.QuoteToGraphic(string(b))
}
