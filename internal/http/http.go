package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moukoublen/goboilerplate/build"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TrafficLogLevel int16

const (
	TrafficLogLevelNone    TrafficLogLevel = 0
	TrafficLogLevelBasic   TrafficLogLevel = 1
	TrafficLogLevelVerbose TrafficLogLevel = 2
)

type Config struct {
	IP                   string
	Port                 int32
	InBoundHTTPLogLevel  TrafficLogLevel
	OutBoundHTTPLogLevel TrafficLogLevel
	LogInLevel           zerolog.Level
	GlobalInboundTimeout time.Duration
	ReadHeaderTimeout    time.Duration
}

// NewDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func NewDefaultRouter(c Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	if c.InBoundHTTPLogLevel > TrafficLogLevelNone {
		router.Use(NewHTTPInboundLoggerMiddleware(c))
	}

	router.Use(middleware.Recoverer)

	if c.GlobalInboundTimeout > 0 {
		router.Use(middleware.Timeout(c.GlobalInboundTimeout))
	}

	router.Get("/about", AboutHandler)

	// for test purposes
	// router.Get("/panic", func(_ http.ResponseWriter, _ *http.Request) { panic("test panic") })

	LogRoutes(router)

	return router
}

func LogRoutes(r *chi.Mux) {
	if log.Logger.GetLevel() > zerolog.DebugLevel {
		return
	}

	routes := []string{}

	walkFunc := func(_ string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		routes = append(routes, route)

		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Error().Err(err).Msg("error during chi walk")
	} else {
		log.Debug().Strs("routes", routes).Msg("http routes")
	}
}

func AboutHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)

	if err := e.Encode(build.GetInfo()); err != nil {
		log.Error().Err(err).Msg("error during json encoding in about handler")
	}
}

// StartListenAndServe creates and runs server.ListenAndServe in a separate go routine.
// Any error produced by ListenAndServe will be sent to fatalErrCh.
// It returns the server struct.
func StartListenAndServe(addr string, handler http.Handler, readHeaderTimeout time.Duration, fatalErrCh chan<- error) *http.Server {
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fatalErrCh <- err
		}
	}()

	return server
}
