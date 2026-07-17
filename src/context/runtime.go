package context

import (
	"os"
	"runtime"
)

// used for caching runtime information
// and testing purposes
var Current *Runtime

type Runtime struct {
	Path       *Path
	Shell      string
	OS         string
	Home       string
	Arch       string
	Hostname   string
	ConfigPath string
}

func Init(shell string) {
	hostname, _ := os.Hostname()

	Current = &Runtime{
		Shell:    shell,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Home:     Home(),
		Hostname: hostname,
		Path:     getPath(),
	}
}

func Home() string {
	if Current != nil {
		return Current.Home
	}

	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}
	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}
