package main

import (
	"context"
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

const fatalErrorsChBufferSize = 10

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cnf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ilog.SetupLog(cnf.Logging)
	log.Info().Msgf("Starting up")

	fatalErrorsCh := make(chan error, fatalErrorsChBufferSize)

	signalCh := channelForSignals(os.Interrupt, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)

	router := ihttp.NewDefaultRouter(cnf.HTTP)
	server := startHTTPServer(cnf.HTTP, router, fatalErrorsCh)

	select {
	case sig := <-signalCh:
		event := log.Warn().Str("signal", sig.String())
		if sigInt, ok := sig.(syscall.Signal); ok {
			event.Int("signalNumber", int(sigInt))
		}
		event.Msg("signal received")
	case fatalErr := <-fatalErrorsCh:
		log.Error().Err(fatalErr).Msg("fatal error received")
		go func() {
			for err := range fatalErrorsCh {
				log.Error().Err(err).Msg("fatal error received")
			}
		}()
	}

	// graceful shut down.
	deadline := time.Now().Add(cnf.ShutdownTimeout)
	dlCtx, dlCancel := context.WithDeadline(ctx, deadline)

	// shutdown other components/services.
	// srv.Shutdown(ctx)
	// shutdown http server.
	if err := server.Shutdown(dlCtx); err != nil {
		log.Warn().Err(err).Msg("error during http server shutdown")
	}

	dlCancel()
	cancel()
	time.Sleep(1 * time.Second)
	close(fatalErrorsCh)

	log.Info().Msgf("Shutdown completed")
}

const signalChannelBufferSize = 10

func channelForSignals(s ...os.Signal) <-chan os.Signal {
	signalCh := make(chan os.Signal, signalChannelBufferSize)
	signal.Notify(signalCh, s...)

	return signalCh
}

func startHTTPServer(config config.HTTP, handler http.Handler, fatalErrCh chan<- error) *http.Server {
	server := ihttp.StartListenAndServe(
		fmt.Sprintf("%s:%d", config.IP, config.Port),
		handler,
		config.ReadHeaderTimeout,
		fatalErrCh,
	)
	log.Info().Msgf("service started at %s:%d", config.IP, config.Port)

	return server
}
