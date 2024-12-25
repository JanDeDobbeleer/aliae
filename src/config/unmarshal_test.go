package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"NoQuotes", "test", "test"},
		{"DoubleQuotes", "\"test\"", "test"},
		{"SingleQuotes", "'test'", "'test'"},
	}

	for _, tc := range tests {
		result := trimQuotes(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestReadDir(t *testing.T) {
	t.Run("ValidDir", func(t *testing.T) {
		testDirPath := filepath.Join("test", "files")
		absPath, _ := filepath.Abs(testDirPath)
		content, err := readDir(absPath)
		assert.NoError(t, err)
		assert.Equal(t, []byte("it exists\nit exists2"), content)
	})

	t.Run("NonExistentDir", func(t *testing.T) {
		_, err := readDir("path/to/nonexistent/dir")
		assert.Error(t, err)
	})
}

func TestRelativePath(t *testing.T) {
	t.Run("RelativePath", func(t *testing.T) {
		absPath, err := filepath.Abs("./test/files")
		assert.NoError(t, err)
		result, err := validatePath(absPath)
		assert.NoError(t, err)
		assert.Equal(t, absPath, result)
	})

	t.Run("Http config", func(t *testing.T) {
		configPathCache = "https://example.com/config.yaml"
		_, err := validatePath("path/to/nonex	istent/dir")
		assert.Error(t, err)
	})
}

func TestIsYamlExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"YamlExtension", "test.yaml", true},
		{"YmlExtension", "test.yml", true},
		{"NoExtension", "test", false},
		{"InvalidExtension", "test.txt", false},
	}

	for _, tc := range tests {
		got := isYAMLExtension(tc.input)
		assert.Equal(t, tc.expected, got, tc.name)
	}
}
