package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/moukoublen/goboilerplate/build"
)

// SetupDefaultRouter returns a *chi.Mux with a default set of middlewares and an "/about" route.
func SetupDefaultRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(120 * time.Second))

	router.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := json.NewEncoder(w)
		_ = e.Encode(build.BuildInfo())
	})

	return router
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
