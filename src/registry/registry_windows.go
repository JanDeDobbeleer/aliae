package registry

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func PersistEnvironmentVariable(name string, value any) {
	// get environment variable and validate if it has the same value
	// if not, persist the new value and exit
	k, err := registry.OpenKey(windows.HKEY_CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return
	}

	defer k.Close()

	_, valType, err := k.GetValue(name, nil)
	if err != nil {
		// if it doesn't exist, persist the new value and exit
		_ = k.SetStringValue(name, value.(string))
		return
	}

	// get string value and match with what we have
	// if it's the same, exit
	// if not, persist the new value and exit
	// if it doesn't exist, persist the new value and exit
	keyValue, _, _ := k.GetStringValue(name)
	if keyValue == value {
		return
	}

	switch valType {
	case windows.REG_SZ:
		_ = k.SetStringValue(name, value.(string))
	case windows.REG_EXPAND_SZ:
		_ = k.SetExpandStringValue(name, value.(string))
	}
}

func PersistPathEntry(path string) {
	k, err := registry.OpenKey(windows.HKEY_CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return
	}

	defer k.Close()

	key := "Path"

	_, _, err = k.GetStringValue(key)
	if err != nil {
		// if it doesn't exist, persist Path with the current path and exit
		_ = k.SetExpandStringValue(key, path)
		return
	}

	keyValue, _, _ := k.GetStringValue(key)

	// split the current value using the ; delimiter and check if we already have this path entry
	entries := strings.Split(keyValue, ";")
	for _, entry := range entries {
		if entry == path {
			return
		}
	}

	keyValue = fmt.Sprintf("%s;%s", path, keyValue)

	_ = k.SetExpandStringValue(key, keyValue)
}
