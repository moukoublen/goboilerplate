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

	"github.com/moukoublen/goboilerplate/internal/config"
	ihttp "github.com/moukoublen/goboilerplate/internal/http"
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

	signalCh := channelForSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	router := ihttp.NewDefaultRouter(cnf.HTTP)
	server, serverErrCh := startHTTPServer(cnf.HTTP, router)

	select {
	case sig := <-signalCh:
		log.Info().Msgf("Signal received: %d %s", sig, sig.String())
	case servErr := <-serverErrCh:
		if !errors.Is(servErr, http.ErrServerClosed) {
			log.Error().Err(servErr).Msg("error returned from http server")
		}
	}

	// graceful shut down
	deadline := time.Now().Add(cnf.ShutdownTimeout)
	dlCtx, dlCancel := context.WithDeadline(ctx, deadline)
	if err := server.Shutdown(dlCtx); err != nil {
		log.Warn().Err(err).Msg("error during http server shutdown")
	}
	// shutdown other components/services
	// srv.Shutdown(ctx)
	dlCancel()
	cancel()
	time.Sleep(1 * time.Second)
	log.Info().Msgf("Shutdown completed")
}

func channelForSignals(s ...os.Signal) <-chan os.Signal {
	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, s...)
	return signalCh
}

func startHTTPServer(config config.HTTP, handler http.Handler) (*http.Server, <-chan error) {
	server, chErr := ihttp.StartListenAndServe(fmt.Sprintf("%s:%d", config.IP, config.Port), handler)
	log.Info().Msgf("service started at %s:%d", config.IP, config.Port)
	return server, chErr
}
