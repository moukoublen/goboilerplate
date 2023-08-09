package config

import "strings"

func buildEnvVarsNamesMapper(levels map[string]any, envVarPrefix string) func(string) string {
	return func(s string) string {
		s = strings.TrimPrefix(s, envVarPrefix)
		s = strings.ToLower(s)
		return buildKey(&strings.Builder{}, ".", "_", levels, strings.Split(s, "_"))
	}
}

func buildKey(b *strings.Builder, levelSep, wordsSep string, levels map[string]any, parts []string) string {
	if len(parts) == 0 {
		return b.String()
	}

	writeWords := func() {
		last := len(parts) - 1
		for i, s := range parts {
			_, _ = b.WriteString(s)
			if i != last {
				_, _ = b.WriteString(wordsSep)
			}
		}
	}

	if len(levels) == 0 {
		writeWords()
		return b.String()
	}

	s := parts[0]

	currentLevel, found := levels[s]
	if !found {
		writeWords()
		return b.String()
	}

	_, _ = b.WriteString(s)
	if len(parts) > 1 {
		_, _ = b.WriteString(levelSep)
	}

	if innerLevels, is := currentLevel.(map[string]any); is {
		return buildKey(b, levelSep, wordsSep, innerLevels, parts[1:])
	}

	return buildKey(b, levelSep, wordsSep, nil, parts[1:])
}
