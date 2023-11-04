package shell

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/process"
)

const (
	UNKNOWN = "unknown"
)

func Name() string {
	pid := os.Getppid()
	p, _ := process.NewProcess(int32(pid))
	name, err := p.Name()
	if err != nil {
		return UNKNOWN
	}

	executable, _ := os.Executable()
	if name == executable {
		p, _ = p.Parent()
		name, err = p.Name()
	}

	if err != nil {
		return UNKNOWN
	}

	return strings.TrimSuffix(name, ".exe")
}
