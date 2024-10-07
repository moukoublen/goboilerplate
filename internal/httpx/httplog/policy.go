package httplog

import (
	"net/http"
	"strings"
)

type HeaderMatcher interface {
	Match(key string, values []string) bool
}

type LogPolicy struct {
	RequestBodyLogPolicy        RequestBodyLogPolicy
	ResponseBodyLogPolicy       ResponseBodyLogPolicy
	ResponseWriterBodyLogPolicy ResponseWriterBodyLogPolicy
	OmitHeaders                 HeaderMatcher
	MaskedValueHeaders          HeaderMatcher
}

func (l LogPolicy) ShouldOmitHeader(key string, values []string) bool {
	if l.OmitHeaders == nil {
		return false
	}

	return l.OmitHeaders.Match(key, values)
}

func (l LogPolicy) ShouldMaskHeader(key string, values []string) bool {
	if l.MaskedValueHeaders == nil {
		return false
	}

	return l.MaskedValueHeaders.Match(key, values)
}

func (l LogPolicy) ShouldLogRequestBody(r *http.Request) bool {
	if l.RequestBodyLogPolicy == nil {
		return DefaultRequestBodyLogPolicy(r)
	}

	return l.RequestBodyLogPolicy(r)
}

func (l LogPolicy) ShouldLogResponseBody(r *http.Response) bool {
	if l.ResponseBodyLogPolicy == nil {
		return DefaultResponseBodyLogPolicy(r)
	}

	return l.ResponseBodyLogPolicy(r)
}

func (l LogPolicy) ShouldLogResponseWriterBody(headers http.Header, statusCode int, body []byte) bool {
	if l.ResponseWriterBodyLogPolicy == nil {
		return DefaultResponseWriterBodyLogPolicy(headers, statusCode, body)
	}

	return l.ResponseWriterBodyLogPolicy(headers, statusCode, body)
}

type RequestBodyLogPolicy func(r *http.Request) bool

var DefaultRequestBodyLogPolicy RequestBodyLogPolicy = func(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	contentEncoding := r.Header.Get("Content-Encoding")
	return len(contentEncoding) == 0 && (strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html"))
}

type ResponseWriterBodyLogPolicy func(headers http.Header, statusCode int, body []byte) bool

var DefaultResponseWriterBodyLogPolicy ResponseWriterBodyLogPolicy = func(headers http.Header, _ int, _ []byte) bool {
	contentType := headers.Get("Content-Type")
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}

type ResponseBodyLogPolicy func(r *http.Response) bool

var DefaultResponseBodyLogPolicy ResponseBodyLogPolicy = func(r *http.Response) bool {
	h := r.Header
	contentType := h.Get("Content-Type")
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}
