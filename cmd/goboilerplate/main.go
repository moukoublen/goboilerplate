package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ifnotnil/daemon"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/moukoublen/goboilerplate/internal/httpx"
	"github.com/moukoublen/goboilerplate/internal/logx"
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
	// pre-init slog with default config
	logger := logx.InitSLog(logx.Config{LogType: logx.LogTypeText, Level: slog.LevelInfo})
	logger.Info("starting up...")

	cnf, err := config.Load(context.Background(), "APP_", defaultConfigs())
	if err != nil {
		logger.Error("error during config init", logx.Error(err))
		os.Exit(1)
	}

	logger = logx.InitSLog(logx.ParseConfig(cnf))

	dmn := daemon.Start(
		context.Background(),
		daemon.WithLogger(logger),
		daemon.WithShutdownGraceDuration(cnf.Duration("shutdown_timeout")),
	)

	httpConf := httpx.ParseConfig(cnf)
	router := httpx.NewDefaultRouter(dmn.CTX(), httpConf, logger)

	// init services / application
	server := httpx.StartListenAndServe(
		fmt.Sprintf("%s:%d", httpConf.IP, httpConf.Port),
		router,
		httpConf.ReadHeaderTimeout,
		dmn.FatalErrorsChannel(),
	)
	logger.InfoContext(dmn.CTX(), "service started", slog.String("bind", fmt.Sprintf("%s:%d", httpConf.IP, httpConf.Port)))

	// set onShutdown for other components/services.
	dmn.OnShutDown(
		func(ctx context.Context) {
			logger.InfoContext(ctx, "shuting down http server")
			if err := server.Shutdown(ctx); err != nil {
				logger.Warn("error during http server shutdown", logx.Error(err))
			}
		},
	)

	dmn.Wait()
}
