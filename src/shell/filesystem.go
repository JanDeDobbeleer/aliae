package shell

import (
	"bytes"
	"os"
	"path/filepath"
	"time"
)

func filesEqual(name string, data []byte) bool {
	existing, err := os.ReadFile(name)
	return err == nil && bytes.Equal(existing, data)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	if filesEqual(path, data) {
		return nil
	}

	// the temp file must be in the same directory as the target,
	// os.Rename is only atomic within the same volume
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}

	defer os.Remove(tmp.Name())

	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}

	// CreateTemp creates the file with 0600
	if err = tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}

// writeFile writes data to path, atomically when possible.
//
// On Windows, replacing a file via rename requires that no process has the
// target open: MoveFileEx(MOVEFILE_REPLACE_EXISTING) fails with
// ERROR_ACCESS_DENIED as long as a single handle exists, even one opened
// with full sharing. Shells load the init script on startup, so brief holds
// are normal — retry those. When the file is still held after the retries
// (a shell keeping it open), fall back to an in-place write: sharing rules
// allow overwriting a file others are reading, we only lose atomicity for
// this write.
func writeFile(path string, data []byte, perm os.FileMode) error {
	const attempts = 4
	wait := 50 * time.Millisecond

	var err error

	for attempt := 1; ; attempt++ {
		if err = writeFileAtomic(path, data, perm); !canRetryWrite(err) || attempt == attempts {
			break
		}

		time.Sleep(wait)
		wait *= 2
	}

	if err == nil || !canRetryWrite(err) {
		return err
	}

	return os.WriteFile(path, data, perm)
}
