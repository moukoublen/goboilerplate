package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/knadh/koanf/v2"
	"github.com/moukoublen/goboilerplate/internal/httpx/httplog"
	"github.com/moukoublen/goboilerplate/internal/logx"
)

type TrafficLogLevel int8

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
		"http.port":                       "8888",
		"http.inbound_traffic_log_level":  defaultInboundTrafficLogLevel,
		"http.outbound_traffic_log_level": defaultOutboundTrafficLogLevel,
	}
}

type Config struct {
	IP                   string
	Port                 int64
	InBoundHTTPLogLevel  TrafficLogLevel
	OutBoundHTTPLogLevel TrafficLogLevel
	GlobalInboundTimeout time.Duration
	ReadHeaderTimeout    time.Duration
}

func ParseConfig(cnf *koanf.Koanf) Config {
	return Config{
		IP:                   cnf.String("http.ip"),
		Port:                 cnf.Int64("http.port"),
		InBoundHTTPLogLevel:  TrafficLogLevel(cnf.Int64("http.inbound_traffic_log_level")),  //nolint:gosec
		OutBoundHTTPLogLevel: TrafficLogLevel(cnf.Int64("http.outbound_traffic_log_level")), //nolint:gosec
		GlobalInboundTimeout: cnf.Duration("http.global_inbound_timeout"),
		ReadHeaderTimeout:    cnf.Duration("http.read_header_timeout"),
	}
}

// NewDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func NewDefaultRouter(ctx context.Context, c Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	switch c.InBoundHTTPLogLevel { //nolint:exhaustive
	case TrafficLogLevelBasic:
		il := httplog.NewHTTPLogger(
			httplog.WithLogInLevel(slog.LevelDebug),
			httplog.WithLogger(logx.GetFromContext(ctx)),
			httplog.OmitBodyFromLog(),
		)
		router.Use(il.Handler)
	case TrafficLogLevelVerbose:
		il := httplog.NewHTTPLogger(
			httplog.WithLogInLevel(slog.LevelDebug),
			httplog.WithLogger(logx.GetFromContext(ctx)),
		)
		router.Use(il.Handler)
	}

	router.Use(middleware.Recoverer)

	if c.GlobalInboundTimeout > 0 {
		router.Use(middleware.Timeout(c.GlobalInboundTimeout))
	}

	router.Get("/about", AboutHandler)
	router.Get("/echo", EchoHandler)
	// for test purposes
	// router.Get("/panic", func(_ http.ResponseWriter, _ *http.Request) { panic("test panic") })

	LogRoutes(ctx, router)

	return router
}

func LogRoutes(ctx context.Context, r *chi.Mux) {
	logger := logx.GetFromContext(ctx)
	if !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	routes := []string{}

	walkFunc := func(_ string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		routes = append(routes, route)

		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		logger.Error("error during chi walk", logx.Error(err))
	} else {
		logger.Debug("http routes", slog.Any("routes", routes))
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
			if !errors.Is(err, http.ErrServerClosed) {
				fatalErrCh <- err
			}
		}
	}()

	return server
}
