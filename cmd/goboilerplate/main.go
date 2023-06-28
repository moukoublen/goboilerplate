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

	cnf, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ilog.SetupLog(parseLogConfig(cnf))
	log.Info().Msgf("Starting up")

	main := internal.NewMain(
		internal.SetShutdownTimeout(cnf.Duration("shutdown_timeout")),
	)

	httpConf := parseHTTPConfig(cnf)
	router := ihttp.NewDefaultRouter(httpConf)

	// init services / application
	server := ihttp.StartListenAndServe(
		fmt.Sprintf("%s:%d", httpConf.IP, httpConf.Port),
		router,
		httpConf.ReadHeaderTimeout,
		main.FatalErrorsChannel(),
	)
	log.Info().Msgf("service started at %s:%d", httpConf.IP, httpConf.Port)

	// set onShutdown for other components/services.
	main.OnShutDown(
		func(ctx context.Context) {
			if err := server.Shutdown(ctx); err != nil {
				log.Warn().Err(err).Msg("error during http server shutdown")
			}
		},
	)

	main.Run(ctx, cancel)
}
