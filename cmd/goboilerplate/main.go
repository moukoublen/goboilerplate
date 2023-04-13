package main

import (
	"context"
	"fmt"

	"github.com/moukoublen/goboilerplate/internal"
	"github.com/moukoublen/goboilerplate/internal/config"
	ihttp "github.com/moukoublen/goboilerplate/internal/http"
	ilog "github.com/moukoublen/goboilerplate/internal/log"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cnf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ilog.SetupLog(cnf.Logging)
	log.Info().Msgf("Starting up")

	orchestrator := internal.NewOrchestrator()

	router := ihttp.NewDefaultRouter(cnf.HTTP)
	server := ihttp.StartListenAndServe(
		fmt.Sprintf("%s:%d", cnf.HTTP.IP, cnf.HTTP.Port),
		router,
		cnf.HTTP.ReadHeaderTimeout,
		orchestrator.FatalErrorsChannel(),
	)
	log.Info().Msgf("service started at %s:%d", cnf.HTTP.IP, cnf.HTTP.Port)

	// set onShutdown for other components/services.
	orchestrator.OnShutDown(func(ctx context.Context) {
		if err := server.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("error during http server shutdown")
		}
	})

	orchestrator.Run(ctx, cancel)
}
