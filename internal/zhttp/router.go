package zhttp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	xhttp "github.com/ifnotnil/x/http"
	"github.com/knadh/koanf/v2"
	"github.com/moukoublen/goboilerplate/internal/zlog"
)

func DefaultConfigValues() map[string]any {
	return map[string]any{
		"http.ip":   "0.0.0.0",
		"http.port": "8888",
	}
}

type Config struct {
	IP                   string
	Port                 int64
	GlobalInboundTimeout time.Duration
	ReadHeaderTimeout    time.Duration
}

func ParseConfig(cnf *koanf.Koanf) Config {
	return Config{
		IP:                   cnf.String("http.ip"),
		Port:                 cnf.Int64("http.port"),
		GlobalInboundTimeout: cnf.Duration("http.global_inbound_timeout"),
		ReadHeaderTimeout:    cnf.Duration("http.read_header_timeout"),
	}
}

// NewDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func NewDefaultRouter(ctx context.Context, c Config, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	router.Use(middleware.Recoverer)

	if c.GlobalInboundTimeout > 0 {
		router.Use(middleware.Timeout(c.GlobalInboundTimeout))
	}

	router.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("", map[string]string{
			"Yoda": "_Named must your fear be before banish it you can_",
		}))
		r.Get("/echo", xhttp.EchoHandler(logger))
	})

	router.Get("/about", AboutHandler)

	// for test purposes
	// router.Get("/panic", func(_ http.ResponseWriter, _ *http.Request) { panic("test panic") })

	LogRoutes(ctx, router)

	return router
}

func LogRoutes(ctx context.Context, r *chi.Mux) {
	logger := zlog.GetFromContext(ctx)
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
		logger.ErrorContext(ctx, "error during chi walk", zlog.Error(err))
	} else {
		logger.DebugContext(ctx, "http routes", slog.Any("routes", routes))
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
