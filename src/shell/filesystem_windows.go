package shell

import (
	"errors"
	"syscall"

	"golang.org/x/sys/windows"
)

// canRetryWrite reports whether the write failed because another process
// holds the file (or its rename target) open, making a retry or an in-place
// write worthwhile. Replacing an open file via rename fails with
// ERROR_ACCESS_DENIED, not a sharing violation: MoveFileEx with
// MOVEFILE_REPLACE_EXISTING requires the target to have no open handles at
// all. All other errors are considered persistent.
func canRetryWrite(err error) bool {
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		return false
	}

	return errno == windows.ERROR_ACCESS_DENIED || errno == windows.ERROR_SHARING_VIOLATION || errno == windows.ERROR_LOCK_VIOLATION
}
