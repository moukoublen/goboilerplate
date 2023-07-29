package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// RespondJSON renders a json response using a json encoder directly over the ResponseWriter.
// That's why in most cases will end up sending chunked (`transfer-encoding: chunked`) response.
func RespondJSON(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Add(`Content-Type`, `application/json; charset=utf-8`)
	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger := zerolog.Ctx(ctx)
			logger.Error().Err(err).Msg("error during response encoding")
		}
	}
}

func NewHTTPInboundLoggerMiddleware(cnf Config) func(http.Handler) http.Handler {
	switch cnf.InBoundHTTPLogLevel {
	case TrafficLogLevelNone:
		return func(h http.Handler) http.Handler { return h }
	case TrafficLogLevelBasic:
		return middleware.RequestLogger(&ChiZerolog{LogInLevel: cnf.LogInLevel})
	case TrafficLogLevelVerbose:
		return RequestResponseLogger(cnf)
	}

	return func(h http.Handler) http.Handler { return h }
}

func RequestResponseLogger(cnf Config) func(next http.Handler) http.Handler {
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
