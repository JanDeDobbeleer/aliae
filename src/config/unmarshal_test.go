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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimQuotes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFile(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		testDirPath := filepath.Join("test", "files", "test.yaml")
		absPath, _ := filepath.Abs(testDirPath)
		content, err := getFile([]string{absPath})
		assert.NoError(t, err)
		assert.Equal(t, "it exists", content)
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := getFile([]string{"path/to/nonexistent/file.txt"})
		assert.Error(t, err)
	})
}

func TestGetDirFiles(t *testing.T) {
	t.Run("ValidDir", func(t *testing.T) {
		testDirPath := filepath.Join("test", "files")
		absPath, _ := filepath.Abs(testDirPath)
		content, err := getDirFiles([]string{absPath})
		assert.NoError(t, err)
		assert.Equal(t, "\nit exists\nit exists2", content)
	})

	t.Run("NonExistentDir", func(t *testing.T) {
		_, err := getDirFiles([]string{"path/to/nonexistent/dir"})
		assert.Error(t, err)
	})
}

func TestRelativePath(t *testing.T) {
	t.Run("RelativePath", func(t *testing.T) {
		absPath, err := filepath.Abs("./test/files")
		assert.NoError(t, err)
		result, err := relativePath(absPath)
		assert.NoError(t, err)
		assert.Equal(t, absPath, result)
	})

	t.Run("HttpConfig", func(t *testing.T) {
		configPathCache = "https://example.com/config.yaml"
		_, err := relativePath("path/to/nonex	istent/dir")
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
