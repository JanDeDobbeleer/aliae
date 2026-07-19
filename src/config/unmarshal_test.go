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
		content, err := readDir(absPath, "")
		assert.NoError(t, err)
		assert.Equal(t, []byte("it exists\nit exists2"), content)
	})

	t.Run("NonExistentDir", func(t *testing.T) {
		content, err := readDir("path/to/nonexistent/dir", "")
		assert.NoError(t, err)
		assert.Equal(t, []byte{}, content)
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
		input := fmt.Sprintf(`env: !include "%s if=eq 1 1"`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "TEST_ENV")
	})

	t.Run("include - false condition skips the file", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include "%s if=eq 1 2"`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("include - false condition never reads a nonexistent file", func(t *testing.T) {
		input := `env: !include "does/not/exist.yaml if=eq 1 2"`
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("include_dir - true condition includes the directory contents", func(t *testing.T) {
		input := fmt.Sprintf(`alias: !include_dir "%s if=eq 1 1"`, aliasesDir)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "test2")
	})

	t.Run("include_dir - false condition never reads the directory", func(t *testing.T) {
		input := `alias: !include_dir "does/not/exist if=eq 1 2"`
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "alias: null", string(result))
	})

	t.Run("list item - false condition drops the entry, leaving the rest intact", func(t *testing.T) {
		input := "alias:\n  - !include \"does/not/exist.yaml if=eq 1 2\"\n  - name: g\n    value: git"
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "alias:\n\n  - name: g\n    value: git", string(result))
	})

	t.Run("path with spaces is preserved alongside a condition", func(t *testing.T) {
		dir := t.TempDir()
		spacedPath := filepath.Join(dir, "my file.yaml")
		assert.NoError(t, os.WriteFile(spacedPath, []byte("- name: sp\n  value: spaced"), 0o644))

		input := fmt.Sprintf(`alias: !include "%s if=eq 1 1"`, spacedPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "spaced")
	})

	t.Run("condition with a quoted argument, single-quoted outer", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include '%s if=eq "a" "a"'`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "TEST_ENV")
	})

	t.Run("condition with a quoted argument, double-quoted outer with escaped quotes", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include "%s if=eq \"a\" \"a\""`, envPath)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "TEST_ENV")
	})

	t.Run("real YAML decode round-trip through a false condition", func(t *testing.T) {
		input := `env: !include "does/not/exist.yaml if=eq 1 2"` + "\n"
		var a Aliae
		err := aliaeUnmarshaler(&a, []byte(input))
		assert.NoError(t, err)
	})
}

func TestExtractSection(t *testing.T) {
	t.Run("bare list is not wrapped", func(t *testing.T) {
		extracted, wrapped, err := extractSection([]byte("- name: test\n  value: test"), "alias")
		assert.NoError(t, err)
		assert.False(t, wrapped)
		assert.Nil(t, extracted)
	})

	t.Run("empty file is not wrapped", func(t *testing.T) {
		extracted, wrapped, err := extractSection([]byte(""), "alias")
		assert.NoError(t, err)
		assert.False(t, wrapped)
		assert.Nil(t, extracted)
	})

	t.Run("map without any recognized section key is not wrapped", func(t *testing.T) {
		extracted, wrapped, err := extractSection([]byte("foo: bar"), "alias")
		assert.NoError(t, err)
		assert.False(t, wrapped)
		assert.Nil(t, extracted)
	})

	t.Run("wrapped file carrying the requested section returns just that section", func(t *testing.T) {
		input := "alias:\n  - name: test\n    value: test\nenv:\n  - name: FOO\n    value: bar\n"
		extracted, wrapped, err := extractSection([]byte(input), "alias")
		assert.NoError(t, err)
		assert.True(t, wrapped)
		assert.Contains(t, string(extracted), "name: test")
		assert.NotContains(t, string(extracted), "FOO")
	})

	t.Run("wrapped file missing the requested section is wrapped with no content", func(t *testing.T) {
		input := "env:\n  - name: FOO\n    value: bar\n"
		extracted, wrapped, err := extractSection([]byte(input), "alias")
		assert.NoError(t, err)
		assert.True(t, wrapped)
		assert.Empty(t, extracted)
	})

	t.Run("nested !include tag inside the extracted section survives untouched", func(t *testing.T) {
		input := "alias:\n  - !include \"nested.yaml\"\n  - name: g\n    value: git\n"
		extracted, wrapped, err := extractSection([]byte(input), "alias")
		assert.NoError(t, err)
		assert.True(t, wrapped)
		assert.Contains(t, string(extracted), `!include "nested.yaml"`)
	})

	t.Run("multiple YAML documents are rejected", func(t *testing.T) {
		input := "alias:\n  - name: a\n    value: a\n---\nenv:\n  - name: B\n    value: b\n"
		_, _, err := extractSection([]byte(input), "alias")
		assert.Error(t, err)
	})
}

func TestIncludeUnmarshalerSection(t *testing.T) {
	t.Run("include - wrapped file contributes only the destination section", func(t *testing.T) {
		dir := t.TempDir()
		combined := filepath.Join(dir, "combined.yaml")
		assert.NoError(t, os.WriteFile(combined, []byte("alias:\n  - name: a\n    value: a\nenv:\n  - name: B\n    value: b\n"), 0o644))

		input := fmt.Sprintf(`alias: !include "%s"`, combined)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "name: a")
		assert.NotContains(t, string(result), "name: B")
	})

	t.Run("include - wrapped file missing the destination section resolves to null", func(t *testing.T) {
		dir := t.TempDir()
		combined := filepath.Join(dir, "combined.yaml")
		assert.NoError(t, os.WriteFile(combined, []byte("env:\n  - name: B\n    value: b\n"), 0o644))

		input := fmt.Sprintf(`alias: !include "%s"`, combined)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "alias: null", string(result))
	})

	t.Run("include - bare list file is unaffected by section extraction", func(t *testing.T) {
		dir := t.TempDir()
		bare := filepath.Join(dir, "bare.yaml")
		assert.NoError(t, os.WriteFile(bare, []byte("- name: a\n  value: a\n"), 0o644))

		input := fmt.Sprintf(`alias: !include "%s"`, bare)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "name: a")
	})

	t.Run("include - list item form is not eligible for section extraction", func(t *testing.T) {
		dir := t.TempDir()
		combined := filepath.Join(dir, "combined.yaml")
		assert.NoError(t, os.WriteFile(combined, []byte("alias:\n  - name: a\n    value: a\n"), 0o644))

		input := fmt.Sprintf("alias:\n  - !include \"%s\"\n", combined)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		// the whole wrapped document is spliced in raw, exactly as before this feature existed
		assert.Contains(t, string(result), "alias:")
		assert.Contains(t, string(result), "name: a")
	})

	t.Run("include_dir - each file contributes only its own destination section", func(t *testing.T) {
		dir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "python.yaml"), []byte("alias:\n  - name: a\n    value: a\nenv:\n  - name: B\n    value: b\n"), 0o644))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "bare.yaml"), []byte("- name: c\n  value: c\n"), 0o644))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "envonly.yaml"), []byte("env:\n  - name: D\n    value: d\n"), 0o644))

		input := fmt.Sprintf(`alias: !include_dir "%s"`, dir)
		result, err := includeUnmarshaler([]byte(input))
		assert.NoError(t, err)
		assert.Contains(t, string(result), "name: a")
		assert.Contains(t, string(result), "name: c")
		assert.NotContains(t, string(result), "name: B")
		assert.NotContains(t, string(result), "name: D")
	})

	t.Run("real decode of a combined workflow file included at two different sections", func(t *testing.T) {
		dir := t.TempDir()
		combined := filepath.Join(dir, "python.yaml")
		assert.NoError(t, os.WriteFile(combined, []byte(
			"alias:\n  - name: a\n    value: a\nenv:\n  - name: B\n    value: b\n",
		), 0o644))

		input := fmt.Sprintf("alias: !include \"%s\"\nenv: !include \"%s\"\n", combined, combined)

		var a Aliae
		err := aliaeUnmarshaler(&a, []byte(input))
		assert.NoError(t, err)
		assert.Len(t, a.Aliae, 1)
		assert.Equal(t, "a", a.Aliae[0].Name)
		assert.Len(t, a.Envs, 1)
		assert.Equal(t, "B", a.Envs[0].Name)
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
