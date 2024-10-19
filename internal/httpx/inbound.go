package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/moukoublen/goboilerplate/internal/logx"
)

// DrainAndCloseRequest can be used (most probably with defer) from the server side to ensure that the http request body is fully consumed and closed.
func DrainAndCloseRequest(r *http.Request, errOut *error) {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return
	}

	_, discardErr := io.Copy(io.Discard, r.Body)
	closeErr := r.Body.Close()

	if discardErr != nil || closeErr != nil {
		*errOut = errors.Join(*errOut, discardErr, closeErr)
	}
}

func ReadJSONRequest(r *http.Request, decodeTo any) (e error) {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return nil
	}

	defer DrainAndCloseRequest(r, &e)

	if err := json.NewDecoder(r.Body).Decode(decodeTo); err != nil {
		return err
	}

	return
}

// RespondJSON renders a json response using a json encoder directly over the ResponseWriter.
// That's why in most cases will end up sending chunked (`transfer-encoding: chunked`) response.
func RespondJSON(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Add(`Content-Type`, `application/json; charset=utf-8`)
	w.Header().Add(`Cache-Control`, `no-store`) // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger := logx.GetFromContext(ctx)
			logger.Error("error during response encoding", logx.Error(err))
		}
	}
}
