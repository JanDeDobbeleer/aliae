//go:build !windows

package shell

import "golang.org/x/sys/unix"

// canReadPath reports whether path is readable by the current user.
func canReadPath(path string) bool {
	return unix.Access(path, unix.R_OK) == nil
}

// canTraverseDir reports whether path is traversable (executable) by the current user.
func canTraverseDir(path string) bool {
	return unix.Access(path, unix.X_OK) == nil
}
