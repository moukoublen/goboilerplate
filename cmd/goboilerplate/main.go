package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	router := ihttp.NewDefaultRouter(cnf.HTTP)

	server := startHTTPServer(cnf.HTTP, router)
	ihttp.LogRoutes(router)

	blockForSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// shutdown the http server gracefully.
	dlCtx, dlCancel := context.WithTimeout(ctx, cnf.ShutdownTimeout)
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

func startHTTPServer(config config.HTTP, handler http.Handler) *http.Server {
	server, chErr := ihttp.StartListenAndServe(fmt.Sprintf("%s:%d", config.IP, config.Port), handler)
	go func() {
		err := <-chErr
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("error returned from http server")
		}
	}()

	log.Info().Msgf("service started at %s:%d", config.IP, config.Port)

	return server
}
