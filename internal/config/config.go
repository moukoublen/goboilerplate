package config

import (
	"time"

	"github.com/rs/zerolog"
)

type HTTPLogLevel int16

const (
	HTTPLogLevelNone            HTTPLogLevel = 0
	HTTPLogLevelBasic           HTTPLogLevel = 1
	HTTPLogLevelRequestResponse HTTPLogLevel = 2
)

type Config struct {
	IP               string
	Port             int32
	ShutdownDeadline time.Duration
	Logging          Logging
}

type Logging struct {
	ConsoleWriter        bool
	LogLevel             zerolog.Level
	InBoundHTTPLogLevel  HTTPLogLevel
	OutBoundHTTPLogLevel HTTPLogLevel
	LogInLevel           zerolog.Level
}

func New() (Config, error) {
	return Config{
		IP:               "0.0.0.0",
		Port:             43000,
		ShutdownDeadline: 4 * time.Second,
		Logging: Logging{
			ConsoleWriter:        false,
			LogLevel:             zerolog.TraceLevel,
			InBoundHTTPLogLevel:  HTTPLogLevelRequestResponse,
			OutBoundHTTPLogLevel: HTTPLogLevelRequestResponse,
			LogInLevel:           zerolog.TraceLevel,
		},
	}, nil
}
