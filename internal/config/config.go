package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

type Config struct {
	HTTP            HTTP          `envPrefix:"HTTP_"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"4s"`
	Logging         Logging       `envPrefix:"LOG_"`
}

type HTTP struct {
	IP                   string        `env:"IP" envDefault:"0.0.0.0"`
	Port                 int32         `env:"PORT" envDefault:"43000"`
	InBoundHTTPLogLevel  int16         `env:"INBOUND_TRAFFIC_LOG_LEVEL" envDefault:"2"`
	OutBoundHTTPLogLevel int16         `env:"OUTBOUND_TRAFFIC_LOG_LEVEL" envDefault:"2"`
	LogInLevel           zerolog.Level `env:"TRAFFIC_LOG_IN_LEVEL" envDefault:"-1"`
	GlobalInboundTimeout time.Duration `env:"GLOBAL_INBOUND_TIMEOUT"`
	ReadHeaderTimeout    time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
}

type Logging struct {
	ConsoleWriter bool          `env:"CONSOLE_WRITER" envDefault:"false"`
	LogLevel      zerolog.Level `env:"LEVEL" envDefault:"-1"`
}

func New() (Config, error) {
	cfg := Config{}
	opts := env.Options{
		Prefix: "APP_",
	}

	if err := env.Parse(&cfg, opts); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
