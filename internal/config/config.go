package config

import (
	"context"
	"errors"
	"io/fs"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/moukoublen/goboilerplate/internal/zlog"
)

func Load(ctx context.Context, envVarPrefix string, defaultConfigs map[string]any) (*koanf.Koanf, error) {
	logger := zlog.GetFromContext(ctx)

	const delim = "."
	k := koanf.New(delim)

	// Load default values.
	if err := k.Load(confmap.Provider(defaultConfigs, delim), nil); err != nil {
		logger.DebugContext(ctx, "error during config loading from defaults", zlog.Error(err))
	}

	// Load YAML config.
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		logger.DebugContext(ctx, "error during config loading from yaml file", zlog.Error(err))
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
		logger.WarnContext(ctx, "error during config loading from env vars", zlog.Error(err))
	}

	// Load .env file
	if err := k.Load(file.Provider(".env"), dotenv.ParserEnv(envVarPrefix, delim, buildEnvVarsNamesMapper(envVarsLevels, envVarPrefix))); err != nil {
		logger.WarnContext(ctx, "error during config loading from dot env file", zlog.Error(err))
		if !errors.Is(err, fs.ErrNotExist) {
			return k, err
		}
	}

	return k, nil
}
