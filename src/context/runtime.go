package context

import (
	"os"
	"runtime"
)

// used for caching runtime information
// and testing purposes
var Current *Runtime

type Runtime struct {
	Shell     string
	OS        string
	Home      string
	ConfigDir string
	CacheDir  string
	Arch      string
	User      string
	Path      *Path
	Hostname  string
}

func Init(shell string) {
	Current = &Runtime{
		Shell:     shell,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		User:      os.Getenv("USER"),
		Home:      Home(),
		ConfigDir: ConfigDir(),
		CacheDir:  CacheDir(),
		Path:      getPath(),
		Hostname:  Hostname(),
	}
}

func ConfigDir() string {
	if Current != nil {
		return Current.ConfigDir
	}
	path, err := os.UserConfigDir()
	if err == nil && len(path) > 0 {
		return path
	}
	return ""
}

func CacheDir() string {
	if Current != nil {
		return Current.CacheDir
	}
	path, err := os.UserCacheDir()
	if err == nil && len(path) > 0 {
		return path
	}
	return ""
}

func Hostname() string {
	if Current != nil {
		return Current.Hostname
	}
	hostname, err := os.Hostname()
	if err == nil && len(hostname) > 0 {
		return hostname
	}
	return ""
}

func Home() string {
	if Current != nil {
		return Current.Home
	}

	home, err := os.UserHomeDir()
	if err == nil && len(home) > 0 {
		return home
	}
	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}
