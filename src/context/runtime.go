package context

import (
	"os"
	"runtime"
)

// used for caching runtime information
// and testing purposes
var Current *Runtime

type Runtime struct {
	Shell string
	OS    string
	Home  string
}

func Init(shell string) {
	Current = &Runtime{
		Shell: shell,
		OS:    runtime.GOOS,
		Home:  Home(),
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
