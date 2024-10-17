package httpx

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/moukoublen/goboilerplate/build"
	"github.com/moukoublen/goboilerplate/internal/logx"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(r.Context(), w, http.StatusOK, build.GetInfo())
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var cErr error
	defer DrainAndCloseRequest(r, &cErr)

	logger := logx.GetFromContext(r.Context())

	if h := r.Header.Get("X-Auth-Key"); h != `_Named must your fear be before banish it you can_` {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b := map[string]any{}
	{
		b["method"] = r.Method
		b["proto"] = r.Proto
		b["protoMajor"] = r.ProtoMajor
		b["protoMinor"] = r.ProtoMinor

		if r.URL != nil {
			u := map[string]any{}
			u["scheme"] = r.URL.Scheme
			u["opaque"] = r.URL.Opaque
			u["user"] = r.URL.User
			u["host"] = r.URL.Host
			u["path"] = r.URL.Path
			u["rawPath"] = r.URL.RawPath
			u["omitHost"] = r.URL.OmitHost
			u["forceQuery"] = r.URL.ForceQuery
			u["rawQuery"] = r.URL.RawQuery
			u["fragment"] = r.URL.Fragment
			u["rawFragment"] = r.URL.RawFragment
		}

		headers := map[string]string{}
		for k := range r.Header {
			headers[k] = r.Header.Get(k)
		}
		b["headers"] = headers

		if ct := r.Header.Get("Content-Type"); strings.Contains(ct, "application/json") {
			body := map[string]any{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			b["body"] = body
		} else {
			body, _ := io.ReadAll(r.Body)
			b["body"] = body
		}

		b["contentLength"] = r.ContentLength
		b["transferEncoding"] = r.TransferEncoding
		b["host"] = r.Host

		trailer := map[string]string{}
		for k := range r.Trailer {
			trailer[k] = r.Header.Get(k)
		}
		b["trailer"] = trailer

		b["remoteAddr"] = r.RemoteAddr
		b["requestURI"] = r.RequestURI
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(b)
	if err != nil {
		logger.Error("error during json encoding", logx.Error(err))
	}
}
