package config

import (
	"errors"
	"io/fs"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HTTP            HTTP
	ShutdownTimeout time.Duration
	Logging         Logging
}

type HTTP struct {
	IP                   string
	Port                 int32
	InBoundHTTPLogLevel  int16
	OutBoundHTTPLogLevel int16
	LogInLevel           zerolog.Level
	GlobalInboundTimeout time.Duration
	ReadHeaderTimeout    time.Duration
}

type Logging struct {
	ConsoleWriter bool
	LogLevel      zerolog.Level
}

func Load(defaultValues map[string]any) (*koanf.Koanf, error) {
	const delim = "."
	k := koanf.New(delim)

	if err := k.Load(confmap.Provider(defaultValues, delim), nil); err != nil {
		log.Warn().Err(err).Msg("error during config loading from defaults")
	}

	// Load JSON config.
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Warn().Err(err).Msg("error during config loading from yaml file")
		if !errors.Is(err, fs.ErrNotExist) {
			return k, err
		}
	}

	// Environment Variables layers
	envVarsLevels := map[string]any{
		"http": map[string]any{},
		"log":  map[string]any{},
	}
	if err := k.Load(env.Provider("APP_", delim, buildEnvVarsNamesMapper(envVarsLevels)), nil); err != nil {
		log.Warn().Err(err).Msg("error during config loading from env vars")
	}

	// Dot env file
	if err := k.Load(file.Provider(".env"), dotenv.ParserEnv("APP_", delim, buildEnvVarsNamesMapper(envVarsLevels))); err != nil {
		log.Warn().Err(err).Msg("error during config loading from dot env file")
		if !errors.Is(err, fs.ErrNotExist) {
			return k, err
		}
	}

	return k, nil
}
