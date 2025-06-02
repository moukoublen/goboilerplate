package zhttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/moukoublen/goboilerplate/internal/zlog"
)

// RespondJSON renders a json response using a json encoder directly over the ResponseWriter.
// That's why in most cases will end up sending chunked (`transfer-encoding: chunked`) response.
func RespondJSON(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Add(`Content-Type`, `application/json; charset=utf-8`)
	w.Header().Add(`Cache-Control`, `no-store`) // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger := zlog.GetFromContext(ctx)
			logger.ErrorContext(ctx, "error during response encoding", zlog.Error(err))
		}
	}
}
