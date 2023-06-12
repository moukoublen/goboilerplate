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

func parseHTTPConfig(c config.HTTP) ihttp.Config {
	return ihttp.Config{
		IP:                   c.IP,
		Port:                 c.Port,
		InBoundHTTPLogLevel:  ihttp.TrafficLogLevel(c.InBoundHTTPLogLevel),
		OutBoundHTTPLogLevel: ihttp.TrafficLogLevel(c.OutBoundHTTPLogLevel),
		LogInLevel:           c.LogInLevel,
		GlobalInboundTimeout: c.GlobalInboundTimeout,
		ReadHeaderTimeout:    c.ReadHeaderTimeout,
	}
}

func parseLogConfig(c config.Logging) ilog.Config {
	return ilog.Config(c)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cnf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ilog.SetupLog(parseLogConfig(cnf.Logging))
	log.Info().Msgf("Starting up")

	main := internal.NewMain(
		internal.SetShutdownTimeout(cnf.ShutdownTimeout),
	)

	router := ihttp.NewDefaultRouter(parseHTTPConfig(cnf.HTTP))

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
