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
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	internal.SetupLog(config.Logging)

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	server := startHTTPServer(config)

	blockForSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// shutdown gracefully the http server.
	dlCtx, dlCancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
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

func startHTTPServer(config config.Config) *http.Server {
	router := internal.SetupDefaultRouter()
	server, chErr := internal.StartListenAndServe(fmt.Sprintf("%s:%d", config.IP, config.Port), router)
	go func() {
		err := <-chErr
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("error returned from http server")
		}
	}()
	internal.LogRoutes(router)
	log.Info().Msgf("service started at %s:%d", config.IP, config.Port)

	return server
}
