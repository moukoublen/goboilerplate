package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moukoublen/goboilerplate/internal"
	"github.com/moukoublen/goboilerplate/internal/config"
	ilog "github.com/moukoublen/goboilerplate/internal/log"
	"github.com/rs/zerolog/log"
)

func main() {
	cnf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	ilog.SetupLog(cnf.Logging)
	log.Info().Msgf("Starting up")

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	router := internal.NewDefaultRouter(cnf.Logging)

	server := startHTTPServer(cnf, router)
	ilog.LogRoutes(router)

	blockForSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// shutdown the http server gracefully.
	dlCtx, dlCancel := context.WithDeadline(ctx, time.Now().Add(4*time.Second))
	_ = server.Shutdown(dlCtx)
	dlCancel()

	cancel()
	log.Info().Msgf("Shutdown completed")
}

func blockForSignals(s ...os.Signal) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, s...)
	sig := <-signalCh
	log.Info().Msgf("Signal received: %d %s\n", sig, sig.String())
	close(signalCh)
}

func startHTTPServer(config config.Config, handler http.Handler) *http.Server {
	server, chErr := internal.StartListenAndServe(fmt.Sprintf("%s:%d", config.IP, config.Port), handler)
	go func() {
		err := <-chErr
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("error returned from http server")
		}
	}()

	log.Info().Msgf("service started at %s:%d", config.IP, config.Port)

	return server
}
