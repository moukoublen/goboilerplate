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

	main := internal.NewMain()

	router := ihttp.NewDefaultRouter(cnf.HTTP)

	// init services / application

	server := ihttp.StartListenAndServe(
		fmt.Sprintf("%s:%d", cnf.HTTP.IP, cnf.HTTP.Port),
		router,
		cnf.HTTP.ReadHeaderTimeout,
		main.FatalErrorsChannel(),
	)
	log.Info().Msgf("service started at %s:%d", cnf.HTTP.IP, cnf.HTTP.Port)

	// set onShutdown for other components/services.
	main.OnShutDown(func(ctx context.Context) {
		if err := server.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("error during http server shutdown")
		}
	})

	main.Run(ctx, cancel)
}
