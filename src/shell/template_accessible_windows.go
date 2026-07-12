//go:build windows

package shell

// canReadPath always returns true on Windows; os.Stat having succeeded is
// treated as sufficient, since Windows ACLs aren't checked with a simple probe.
func canReadPath(_ string) bool {
	return true
}

// canTraverseDir always returns true on Windows; os.Stat having succeeded is
// treated as sufficient, since Windows ACLs aren't checked with a simple probe.
func canTraverseDir(_ string) bool {
	return true
}
