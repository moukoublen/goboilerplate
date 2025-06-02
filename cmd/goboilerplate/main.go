package main

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"

	"github.com/ifnotnil/daemon"
	"github.com/moukoublen/goboilerplate/internal/config"
	"github.com/moukoublen/goboilerplate/internal/zhttp"
	"github.com/moukoublen/goboilerplate/internal/zlog"
)

func defaultConfigs() map[string]any {
	defaults := map[string]any{
		"shutdown_timeout": "4s",
	}

	gather := []map[string]any{
		zhttp.DefaultConfigValues(),
		zlog.DefaultConfigValues(),
	}

	for _, g := range gather {
		maps.Copy(defaults, g)
	}

	return defaults
}

func main() {
	// pre-init slog with default config
	logger := zlog.InitSLog(zlog.Config{LogType: zlog.LogTypeText, Level: slog.LevelInfo})
	logger.Info("starting up...")

	cnf, err := config.Load(context.Background(), "APP_", defaultConfigs())
	if err != nil {
		logger.Error("error during config init", zlog.Error(err))
		os.Exit(1)
	}

	logger = zlog.InitSLog(zlog.ParseConfig(cnf))

	dmn := daemon.Start(
		context.Background(),
		daemon.WithLogger(logger),
		daemon.WithShutdownGraceDuration(cnf.Duration("shutdown_timeout")),
	)

	httpConf := zhttp.ParseConfig(cnf)
	router := zhttp.NewDefaultRouter(dmn.CTX(), httpConf, logger)

	// init services / application
	server := zhttp.StartListenAndServe(
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
				logger.Warn("error during http server shutdown", zlog.Error(err))
			}
		},
	)

	dmn.Wait()
}
