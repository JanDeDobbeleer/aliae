package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestResolveConfigPath(t *testing.T) {
	cases := []struct {
		name      string
		configVar string
		homeVar   string
		expected  string
	}{
		{"Config env var", "test", "", "test"},
		{"No config env var", "", "/home", "/home/.aliae.yaml"},
		{"No config env var, no home", "", "", ".aliae.yaml"},
	}

	for _, c := range cases {
		os.Setenv("ALIAE_CONFIG", c.configVar)
		os.Setenv("HOME", c.homeVar)
		got := resolveConfigPath("")
		assert.Equal(t, got, got, c.name)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expected    *Aliae
		expectError bool
	}{
		{
			"Valid",
			"aliae.valid.yaml",
			&Aliae{
				Aliae: shell.Aliae{
					{Name: "test", Value: shell.Template("test")},
					{Name: "test2", Value: shell.Template("test2")},
				},
				Envs: shell.Envs{
					{Name: "TEST_ENV", Value: "test"},
				},
			},
			false,
		},
		{
			"Invalid",
			"aliae.invalid.yaml",
			nil,
			true,
		},
	}

	for _, tc := range tests {
		configFile := filepath.Join("test", tc.config)
		got, err := LoadConfig(configFile)

		if tc.expectError {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}

		assert.Equal(t, tc.expected, got, tc.name)
	}
}
