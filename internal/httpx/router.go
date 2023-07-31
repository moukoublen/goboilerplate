package httpx

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/knadh/koanf/v2"
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

// Default configuration values.
const (
	defaultInboundTrafficLogLevel  = 2
	defaultOutboundTrafficLogLevel = 2
)

func DefaultConfigValues() map[string]any {
	return map[string]any{
		"http.ip":                         "0.0.0.0",
		"http.port":                       "43000",
		"http.inbound_traffic_log_level":  defaultInboundTrafficLogLevel,
		"http.outbound_traffic_log_level": defaultOutboundTrafficLogLevel,
		"http.log_in_level":               0,
	}
}

type Config struct {
	IP                   string
	Port                 int32
	InBoundHTTPLogLevel  TrafficLogLevel
	OutBoundHTTPLogLevel TrafficLogLevel
	LogInLevel           zerolog.Level
	GlobalInboundTimeout time.Duration
	ReadHeaderTimeout    time.Duration
}

func ParseConfig(cnf *koanf.Koanf) Config {
	return Config{
		IP:                   cnf.String("http.ip"),
		Port:                 int32(cnf.Int64("http.port")),
		InBoundHTTPLogLevel:  TrafficLogLevel(cnf.Int64("http.inbound_traffic_log_level")),
		OutBoundHTTPLogLevel: TrafficLogLevel(cnf.Int64("http.outbound_traffic_log_level")),
		LogInLevel:           zerolog.Level(cnf.Int64("http.log_in_level")),
		GlobalInboundTimeout: cnf.Duration("http.global_inbound_timeout"),
		ReadHeaderTimeout:    cnf.Duration("http.read_header_timeout"),
	}
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

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(r.Context(), w, http.StatusOK, build.GetInfo())
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
