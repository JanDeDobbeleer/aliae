//go:build !windows

package registry

func PersistEnvironmentVariable(_ string, _ any) {}

func PersistPathEntry(_ string) {}
