package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/moukoublen/goboilerplate/build"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func SetupDefaultRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestLogger(&chiZerolog{}))
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(120 * time.Second))

	router.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := json.NewEncoder(w)
		_ = e.Encode(build.BuildInfo())
	})

	// for test purposes
	router.Get("/panic", AboutRouteHandler)

	return router
}

func AboutRouteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	_ = e.Encode(build.BuildInfo())
}

// StartListenAndServe creates and runs server.ListenAndServe in a separate go routine.
// It returns the server struct and a channel of errors in which will be forewarded any error returned from server.ListenAndServe.
func StartListenAndServe(addr string, router *chi.Mux) (*http.Server, <-chan error) {
	server := &http.Server{Addr: addr, Handler: router}
	errorChannel := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	return server, errorChannel
}

// Log functionality

func LogRoutes(r *chi.Mux) {
	if log.Logger.GetLevel() > zerolog.DebugLevel {
		return
	}

	var traverse func(routes []chi.Route)
	traverse = func(routes []chi.Route) {
		for _, c := range routes {
			for k := range c.Handlers {
				log.Debug().Msgf("Route: %s %s", k, c.Pattern)
			}
			if c.SubRoutes != nil {
				traverse(c.SubRoutes.Routes())
			}
		}
	}
	traverse(r.Routes())
}

type chiZerolog struct{}

func (c *chiZerolog) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &chiZerologEntry{
		request: r,
	}

	logEvent := log.Debug()

	reqID := middleware.GetReqID(r.Context())
	logEvent.Str("request_id", reqID)
	logEvent.Str("method", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logEvent.Str("url", fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto))
	logEvent.Str("from", r.RemoteAddr)

	entry.logEvent = logEvent

	return entry
}

type chiZerologEntry struct {
	request  *http.Request
	logEvent *zerolog.Event
}

func (e *chiZerologEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.logEvent.Int("status", status).Int("bytes", bytes).Dur("elapsed", elapsed).Interface("extra", extra).Msg("http call")
}

func (e *chiZerologEntry) Panic(v interface{}, stack []byte) {
	log.Error().Stack().Interface("item", v).Strs("stack", strings.Split(string(stack), "\n")).Msg("panic during incoming http call")
}
