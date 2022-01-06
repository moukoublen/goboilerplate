package config

import "github.com/rs/zerolog"

type Config struct {
	IP      string
	Port    int32
	Logging Logging
	HTTP    HTTP
}

type Logging struct {
	ConsoleLog bool
	LogLevel   zerolog.Level
}

type HTTP struct {
	EnableLogger bool
	LogInLevel   zerolog.Level
}

func New() (Config, error) {
	return Config{
		IP:   "0.0.0.0",
		Port: 43000,
		Logging: Logging{
			ConsoleLog: true,
			LogLevel:   zerolog.TraceLevel,
		},
		HTTP: HTTP{
			EnableLogger: true,
			LogInLevel:   zerolog.TraceLevel,
		},
	}, nil
}
