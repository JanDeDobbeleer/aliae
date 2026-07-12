package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFile(t *testing.T) {
	cases := []struct {
		Case     string
		Existing string
		Data     string
	}{
		{
			Case: "new file",
			Data: "alias foo = bar",
		},
		{
			Case:     "unchanged content",
			Existing: "alias foo = bar",
			Data:     "alias foo = bar",
		},
		{
			Case:     "changed content",
			Existing: "alias foo = bar",
			Data:     "alias foo = baz",
		},
	}

	for _, tc := range cases {
		path := filepath.Join(t.TempDir(), "aliae.nu")

		if len(tc.Existing) != 0 {
			require.NoError(t, os.WriteFile(path, []byte(tc.Existing), 0o644), tc.Case)
		}

		err := writeFile(path, []byte(tc.Data), 0o644)
		require.NoError(t, err, tc.Case)

		got, err := os.ReadFile(path)
		require.NoError(t, err, tc.Case)
		assert.Equal(t, tc.Data, string(got), tc.Case)

		// no temp files left behind
		entries, err := os.ReadDir(filepath.Dir(path))
		require.NoError(t, err, tc.Case)
		assert.Len(t, entries, 1, tc.Case)
	}
}
