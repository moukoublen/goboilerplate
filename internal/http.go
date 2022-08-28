package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moukoublen/goboilerplate/build"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/moukoublen/goboilerplate/internal/log"
)

// NewDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func NewDefaultRouter(c config.Logging) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	if c.InBoundHTTPLogLevel > config.HTTPLogLevelNone {
		router.Use(log.NewHTTPInboundLoggerMiddleware(c))
	}
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(120 * time.Second))

	router.Get("/about", AboutHandler)

	// for test purposes
	router.Get("/panic", Panic)

	return router
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
		err := server.ListenAndServe()
		if err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	return server, errorChannel
}
