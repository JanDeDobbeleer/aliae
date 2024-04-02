package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	text := `{{ formatString .Value}}`
	cases := []struct {
		Case     string
		Value    interface{}
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
		Value    interface{}
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
		got := ""
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
		Value    interface{}
		Expected string
	}{
		{
			Case:     "string",
			Value:    `hello`,
			Expected: `hello`,
		},
		{
			Case:     "stringWithQuotes",
			Value:    `hello "world"`,
			Expected: `hello \"world\"`,
		},
		{
			Case:     "stringWithBackslashes",
			Value:    `hello \world`,
			Expected: `hello \\world`,
		},
		{
			Case:     "template",
			Value:    Template(`hello "world"`),
			Expected: `hello \"world\"`,
		},
	}

	for _, tc := range cases {
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
