package config

import (
	"fmt"
	"os"
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

func TestIncludeUnmarshalerCondition(t *testing.T) {
	envPath, err := filepath.Abs(filepath.Join("test", "envs", "env.yaml"))
	assert.NoError(t, err)

	aliasesDir, err := filepath.Abs(filepath.Join("test", "aliases"))
	assert.NoError(t, err)

	t.Run("include - true condition includes the file", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include %q if="eq 1 1"`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "TEST_ENV")
	})

	t.Run("include - false condition skips the file", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include %q if="eq 1 2"`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("include - false condition never reads a nonexistent file", func(t *testing.T) {
		input := `env: !include "does/not/exist.yaml" if="eq 1 2"`
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("include_dir - true condition includes the directory contents", func(t *testing.T) {
		input := fmt.Sprintf(`alias: !include_dir %q if="eq 1 1"`, aliasesDir)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "test2")
	})

	t.Run("include_dir - false condition never reads the directory", func(t *testing.T) {
		input := `alias: !include_dir "does/not/exist" if="eq 1 2"`
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "alias: null", string(result))
	})

	t.Run("list item - false condition drops the entry, leaving the rest intact", func(t *testing.T) {
		input := "alias:\n  - !include \"does/not/exist.yaml\" if=\"eq 1 2\"\n  - name: g\n    value: git"
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "alias:\n\n  - name: g\n    value: git", string(result))
	})

	t.Run("path with spaces is preserved alongside a condition", func(t *testing.T) {
		dir := t.TempDir()
		spacedPath := filepath.Join(dir, "my file.yaml")
		assert.NoError(t, os.WriteFile(spacedPath, []byte("- name: sp\n  value: spaced"), 0o644))

		input := fmt.Sprintf(`alias: !include "%s" if="eq 1 1"`, spacedPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "spaced")
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
