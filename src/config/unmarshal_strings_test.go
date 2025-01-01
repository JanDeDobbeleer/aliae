package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "standard string", input: "test", expected: "test"},
		{
			name:     "folded block scalar",
			expected: "one two three",
			input: `>
      one
	  two
	  three`,
		},
		{
			name:     "folded block scalar, multiple >",
			expected: "one two > three",
			input: `>
      one
	  two >
	  three`,
		},
		{
			name: "literal block scalar",
			expected: `one
two`,
			input: `|
      one
      two`,
		},
	}

	for _, tc := range tests {
		var result string
		_ = stringUmarshaler(&result, []byte(tc.input))
		assert.Equal(t, tc.expected, result)
	}
}
