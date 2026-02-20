package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	text := `{{ formatString .Value}}`
	cases := []struct {
		Case     string
		Value    any
		Expected string
	}{
		{
			Case:     "string",
			Value:    "hello",
			Expected: `"hello"`,
		},
		{
			Case:     "bool",
			Value:    true,
			Expected: `true`,
		},
		{
			Case:     "int",
			Value:    32,
			Expected: `32`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: BASH}
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

// This tests both formatArray and splitString
func TestFormatArray(t *testing.T) {
	text := `{{ formatArray .Value }}`
	textDelim := `{{ formatArray .Value .Delim }}`
	cases := []struct {
		Case     string
		Value    any
		Expected string
		Delim    string
	}{
		{
			Case:     "string",
			Value:    "hello",
			Expected: `"hello"`,
		},
		{
			Case:     "Multiple Strings",
			Value:    "hello world, I am a long string",
			Expected: `"hello" "world," "I" "am" "a" "long" "string"`,
		},
		{
			Case: "Multiline String",
			Value: `hello
world
I
am
a
multiline
string`,
			Expected: `"hello" "world" "I" "am" "a" "multiline" "string"`,
		},
		{
			Case: "Single Line Starts with newline",
			Value: `
hello world I am a long string`,
			Expected: `"hello" "world" "I" "am" "a" "long" "string"`,
		},
		{
			Case:     "Single line with delimiter",
			Value:    `hello world I am a long string`,
			Delim:    ",",
			Expected: `"hello","world","I","am","a","long","string"`,
		},
		{
			Case: "Multiline with delimiter",
			Value: `hello
I
am
a
mutliline
string`,
			Delim:    ";",
			Expected: `"hello";"I";"am";"a";"mutliline";"string"`,
		},
		{
			Case:     "bool",
			Value:    true,
			Expected: `true`,
		},
		{
			Case:     "int",
			Value:    32,
			Expected: `32`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: BASH}
		var got string
		if tc.Delim == "" {
			got, _ = parse(text, tc)
		} else {
			got, _ = parse(textDelim, tc)
		}
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestEscapeString(t *testing.T) {
	text := `{{ escapeString .Value}}`
	cases := []struct {
		Case     string
		Shell    string
		Value    any
		Expected string
	}{
		{
			Case:     "string",
			Value:    `hello`,
			Expected: `hello`,
		},
		{
			Case:     "string with quotes",
			Value:    `hello "world"`,
			Expected: `hello \"world\"`,
		},
		{
			Case:     "string with backslashes",
			Value:    `hello \world`,
			Expected: `hello \\world`,
		},
		{
			Case:     "template",
			Value:    Template(`hello "world"`),
			Expected: `hello \"world\"`,
		},
		{
			Case:     "PowerShell: string",
			Shell:    PWSH,
			Value:    `hello`,
			Expected: `hello`,
		},
		{
			Case:     "PowerShell: string with quotes",
			Shell:    PWSH,
			Value:    `hello "world"`,
			Expected: "hello `\"world`\"",
		},
		{
			Case:     "PowerShell: string with backticks",
			Shell:    PWSH,
			Value:    "hello `world",
			Expected: "hello ``world",
		},
		{
			Case:     "PowerShell: template",
			Shell:    PWSH,
			Value:    Template(`hello "world"`),
			Expected: "hello `\"world`\"",
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell}
		if len(tc.Shell) == 0 {
			tc.Shell = BASH
		}
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestMatch(t *testing.T) {
	text := `{{ match .Variable "hello" "world"}}`
	cases := []struct {
		Case     string
		Variable string
		Expected string
	}{
		{
			Case:     "match",
			Variable: "hello",
			Expected: `true`,
		},
		{
			Case:     "match",
			Variable: "world",
			Expected: `true`,
		},
		{
			Case:     "noMatch",
			Variable: "goodbye",
			Expected: `false`,
		},
	}

	for _, tc := range cases {
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestHasCommand(t *testing.T) {
	text := `{{ hasCommand .Command}}`
	cases := []struct {
		Case     string
		Command  string
		Expected string
	}{
		{
			Case:     "hasCommand",
			Command:  "go",
			Expected: `true`,
		},
		{
			Case:     "noCommand",
			Command:  "notACommand",
			Expected: `false`,
		},
	}

	for _, tc := range cases {
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestHomeFileExists(t *testing.T) {
	text := `{{ homeFileExists .Path }}`
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: BASH, Home: tempDir}
	relExisting := filepath.Join(tempDir, ".cache", "aliae")
	absExisting := filepath.Join(tempDir, "absolute.txt")
	cached := filepath.Join(tempDir, "cached.txt")
	assert.NoError(t, os.MkdirAll(filepath.Dir(relExisting), 0o700))
	assert.NoError(t, os.WriteFile(relExisting, []byte("ok"), 0o600))
	assert.NoError(t, os.WriteFile(absExisting, []byte("ok"), 0o600))

	t.Cleanup(clearPathExistsCache)

	cases := []struct {
		Case     string
		Path     string
		Expected string
	}{
		{
			Case:     "relative to home file exists",
			Path:     ".cache/aliae",
			Expected: "true",
		},
		{
			Case:     "absolute file exists",
			Path:     absExisting,
			Expected: "true",
		},
		{
			Case:     "relative file does not exist",
			Path:     ".cache/missing",
			Expected: "false",
		},
	}

	for _, tc := range cases {
		clearPathExistsCache()
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}

	clearPathExistsCache()
	initial, _ := parse(text, struct{ Path string }{Path: "cached.txt"})
	assert.Equal(t, "false", initial, "cached non-existent file")
	assert.NoError(t, os.WriteFile(cached, []byte("now exists"), 0o600))
	cachedResult, _ := parse(text, struct{ Path string }{Path: "cached.txt"})
	assert.Equal(t, "false", cachedResult, "should reuse cached result")
}

func TestHomeDirExists(t *testing.T) {
	text := `{{ homeDirExists .Path }}`
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: BASH, Home: tempDir}
	relDir := filepath.Join(tempDir, ".cache", "aliae")
	absDir := filepath.Join(tempDir, "absolute-dir")
	assert.NoError(t, os.MkdirAll(relDir, 0o700))
	assert.NoError(t, os.MkdirAll(absDir, 0o700))

	t.Cleanup(clearPathExistsCache)

	cases := []struct {
		Case     string
		Path     string
		Expected string
	}{
		{
			Case:     "relative to home directory exists",
			Path:     ".cache/aliae",
			Expected: "true",
		},
		{
			Case:     "absolute directory exists",
			Path:     absDir,
			Expected: "true",
		},
		{
			Case:     "directory does not exist",
			Path:     ".cache/missing-dir",
			Expected: "false",
		},
	}

	for _, tc := range cases {
		clearPathExistsCache()
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}
