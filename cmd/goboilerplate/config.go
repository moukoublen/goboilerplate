package main

import (
	"github.com/knadh/koanf/v2"
	ihttp "github.com/moukoublen/goboilerplate/internal/http"
	ilog "github.com/moukoublen/goboilerplate/internal/log"
	"github.com/rs/zerolog"
)

func parseHTTPConfig(cnf *koanf.Koanf) ihttp.Config {
	return ihttp.Config{
		IP:                   cnf.String("http.ip"),
		Port:                 int32(cnf.Int64("http.port")),
		InBoundHTTPLogLevel:  ihttp.TrafficLogLevel(cnf.Int64("http.inbound_traffic_log_level")),
		OutBoundHTTPLogLevel: ihttp.TrafficLogLevel(cnf.Int64("http.outbound_traffic_log_level")),
		LogInLevel:           zerolog.Level(cnf.Int64("http.log_in_level")),
		GlobalInboundTimeout: cnf.Duration("http.global_inbound_timeout"),
		ReadHeaderTimeout:    cnf.Duration("http.read_header_timeout"),
	}
}

func parseLogConfig(cnf *koanf.Koanf) ilog.Config {
	return ilog.Config{
		ConsoleWriter: cnf.Bool("log.console_writer"),
		LogLevel:      zerolog.Level(cnf.Int("log.level")),
	}
}
