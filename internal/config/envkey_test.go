package config

import (
	"strings"
	"testing"
)

func TestBuildKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		levels   map[string]any
		parts    []string
		expected string
	}{
		"nil": {
			levels:   map[string]any{},
			parts:    []string{},
			expected: "",
		},

		"level 0": {
			levels:   map[string]any{},
			parts:    []string{"a", "b", "c"},
			expected: "a_b_c",
		},

		"level 1": {
			levels: map[string]any{
				"a": map[string]any{},
			},
			parts:    []string{"a", "b", "c"},
			expected: "a.b_c",
		},

		"level 2": {
			levels: map[string]any{
				"a": map[string]any{
					"b": map[string]any{},
				},
			},
			parts:    []string{"a", "b", "c", "d"},
			expected: "a.b.c_d",
		},

		"level 3 edge": {
			levels: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{},
					},
				},
			},
			parts:    []string{"a", "b", "c"},
			expected: "a.b.c",
		},

		"level 3 not found": {
			levels: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{},
					},
				},
			},
			parts:    []string{"a", "b", "d"},
			expected: "a.b.d",
		},

		"level 3 with nil": {
			levels: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": nil,
					},
				},
			},
			parts:    []string{"a", "b", "c", "d"},
			expected: "a.b.c.d",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := buildKey(&strings.Builder{}, ".", "_", tc.levels, tc.parts)
			if got != tc.expected {
				t.Errorf("expected %s got %s", tc.expected, got)
			}
		})
	}
}
