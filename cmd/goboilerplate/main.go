package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/moukoublen/goboilerplate/internal"
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

	fatalBufferSize := 4
	signalsCh := internal.ChannelForSignals(fatalBufferSize, []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT})
	fatalErrorsCh := make(chan error, fatalBufferSize)

	ctx, cancel := context.WithCancel(context.Background())
	router := ihttp.NewDefaultRouter(cnf.HTTP)

	// components/services initialization

	server := ihttp.StartListenAndServe(fmt.Sprintf("%s:%d", cnf.HTTP.IP, cnf.HTTP.Port), router, cnf.HTTP.ReadHeaderTimeout, fatalErrorsCh)
	log.Info().Msgf("service started at %s:%d", cnf.HTTP.IP, cnf.HTTP.Port)

	internal.WaitForFatal(signalsCh, fatalErrorsCh)

	// graceful shut down
	deadline := time.Now().Add(cnf.ShutdownTimeout)
	dlCtx, dlCancel := context.WithDeadline(ctx, deadline)
	if err := server.Shutdown(dlCtx); err != nil {
		log.Warn().Err(err).Msg("error during http server shutdown")
	}

	// shutdown components/services
	// e.g. srv.Shutdown(dlCtx)

	dlCancel()
	cancel()
	time.Sleep(1 * time.Second)
	log.Info().Msgf("Shutdown completed")
}
