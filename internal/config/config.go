package config

import "github.com/rs/zerolog"

type Config struct {
	IP      string
	Port    int32
	Logging Logging
}

type Logging struct {
	ConsoleLog bool
	LogLevel   zerolog.Level
}

func New() (Config, error) {
	return Config{
		IP:   "0.0.0.0",
		Port: 43000,
		Logging: Logging{
			ConsoleLog: true,
			LogLevel:   zerolog.DebugLevel,
		},
	}, nil
}
