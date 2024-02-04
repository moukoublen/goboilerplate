package main

import (
	"context"
	"fmt"

	"github.com/moukoublen/goboilerplate/internal"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/moukoublen/goboilerplate/internal/httpx"
	"github.com/moukoublen/goboilerplate/internal/logx"
	"github.com/rs/zerolog/log"
)

func defaultConfigs() map[string]any {
	defaults := map[string]any{
		"shutdown_timeout": "4s",
	}

	gather := []map[string]any{
		httpx.DefaultConfigValues(),
		logx.DefaultConfigValues(),
	}

	for _, g := range gather {
		for k, v := range g {
			defaults[k] = v
		}
	}

	return defaults
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cnf, err := config.Load("APP_", defaultConfigs())
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	logx.SetupLog(logx.ParseConfig(cnf))
	log.Info().Msgf("Starting up")

	daemon := internal.NewDaemon(
		internal.SetShutdownTimeout(cnf.Duration("shutdown_timeout")),
	)

	httpConf := httpx.ParseConfig(cnf)
	router := httpx.NewDefaultRouter(httpConf)

	// init services / application
	server := httpx.StartListenAndServe(
		fmt.Sprintf("%s:%d", httpConf.IP, httpConf.Port),
		router,
		httpConf.ReadHeaderTimeout,
		daemon.FatalErrorsChannel(),
	)
	log.Info().Msgf("service started at %s:%d", httpConf.IP, httpConf.Port)

	// set onShutdown for other components/services.
	daemon.OnShutDown(
		func(ctx context.Context) {
			if err := server.Shutdown(ctx); err != nil {
				log.Warn().Err(err).Msg("error during http server shutdown")
			}
		},
	)

	daemon.Run(ctx, cancel)
}
