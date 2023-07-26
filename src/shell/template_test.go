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
