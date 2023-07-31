package config

import (
	"errors"
	"io/fs"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func Load(envVarPrefix string) (*koanf.Koanf, error) {
	const delim = "."
	k := koanf.New(delim)

	if err := k.Load(confmap.Provider(defaultConfigs(), delim), nil); err != nil {
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
	if err := k.Load(env.Provider(envVarPrefix, delim, buildEnvVarsNamesMapper(envVarsLevels, envVarPrefix)), nil); err != nil {
		log.Warn().Err(err).Msg("error during config loading from env vars")
	}

	// Dot env file
	if err := k.Load(file.Provider(".env"), dotenv.ParserEnv(envVarPrefix, delim, buildEnvVarsNamesMapper(envVarsLevels, envVarPrefix))); err != nil {
		log.Warn().Err(err).Msg("error during config loading from dot env file")
		if !errors.Is(err, fs.ErrNotExist) {
			return k, err
		}
	}

	return k, nil
}

// Default configuration values.
const (
	defaultInboundTrafficLogLevel  = 2
	defaultOutboundTrafficLogLevel = 2
)

func defaultConfigs() map[string]any {
	return map[string]any{
		"shutdown_timeout":                "4s",
		"http.ip":                         "0.0.0.0",
		"http.port":                       "43000",
		"http.inbound_traffic_log_level":  defaultInboundTrafficLogLevel,
		"http.outbound_traffic_log_level": defaultOutboundTrafficLogLevel,
		"http.log_in_level":               -1,
		"http.read_header_timeout":        "5s",
		"log.console_writer":              false,
		"log.level":                       -1,
	}
}
