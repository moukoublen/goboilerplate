package main

import (
	"github.com/knadh/koanf/v2"
	"github.com/moukoublen/goboilerplate/internal/httpx"
	"github.com/moukoublen/goboilerplate/internal/logx"
	"github.com/rs/zerolog"
)

func parseHTTPConfig(cnf *koanf.Koanf) httpx.Config {
	return httpx.Config{
		IP:                   cnf.String("http.ip"),
		Port:                 int32(cnf.Int64("http.port")),
		InBoundHTTPLogLevel:  httpx.TrafficLogLevel(cnf.Int64("http.inbound_traffic_log_level")),
		OutBoundHTTPLogLevel: httpx.TrafficLogLevel(cnf.Int64("http.outbound_traffic_log_level")),
		LogInLevel:           zerolog.Level(cnf.Int64("http.log_in_level")),
		GlobalInboundTimeout: cnf.Duration("http.global_inbound_timeout"),
		ReadHeaderTimeout:    cnf.Duration("http.read_header_timeout"),
	}
}

func parseLogConfig(cnf *koanf.Koanf) logx.Config {
	return logx.Config{
		ConsoleWriter: cnf.Bool("log.console_writer"),
		LogLevel:      zerolog.Level(cnf.Int("log.level")),
	}
}
