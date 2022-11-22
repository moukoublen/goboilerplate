package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moukoublen/goboilerplate/build"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct{}

// NewDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func NewDefaultRouter(c config.HTTP) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	if c.InBoundHTTPLogLevel > config.HTTPTrafficLogLevelNone {
		router.Use(NewHTTPInboundLoggerMiddleware(c))
	}
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(120 * time.Second))

	router.Get("/about", AboutHandler)

	// for test purposes
	router.Get("/panic", Panic)

	LogRoutes(router)
	return router
}

func LogRoutes(r *chi.Mux) {
	if log.Logger.GetLevel() > zerolog.DebugLevel {
		return
	}

	routes := make([]string, 0, 10)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		routes = append(routes, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Error().Err(err).Msg("error during chi walk")
	} else {
		log.Debug().Strs("routes", routes).Msg("error during chi walk")
	}
}

func Panic(w http.ResponseWriter, r *http.Request) { panic("test panic") }

func AboutHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	_ = e.Encode(build.BuildInfo())
}

// StartListenAndServe creates and runs server.ListenAndServe in a separate go routine.
// It returns the server struct and a channel of errors in which will be forwarded any error returned from server.ListenAndServe.
func StartListenAndServe(addr string, handler http.Handler) (*http.Server, <-chan error) {
	server := &http.Server{Addr: addr, Handler: handler}
	errorChannel := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()
	return server, errorChannel
}
