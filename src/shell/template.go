package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"text/template"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Template string

var (
	pathExistsCache   = map[string]pathInfo{}
	pathExistsCacheMu sync.RWMutex
)

type pathInfo struct {
	exists bool
	isDir  bool
}

func (t Template) Parse() Template {
	value, err := parse(string(t), context.Current)
	if err != nil {
		return t
	}

	return Template(value)
}

func (t Template) String() string {
	return string(t.Parse())
}

func parse(text string, ctx any) (string, error) {
	if !strings.Contains(text, "{{") || !strings.Contains(text, "}}") {
		return text, nil
	}

	parsedTemplate, err := template.New("alias").Funcs(funcMap()).Parse(text)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = parsedTemplate.Execute(buffer, ctx)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"isPwshOption": isPwshOption,
		"isPwshScope":  isPwshScope,
		"formatString": formatString,
		"formatArray":  formatArray,
		"escapeString": escapeString,
		"env":          os.Getenv,
		"match":        match,
		"hasCommand":   hasCommand,
		"homeFileExists": homeFileExists,
		"homeDirExists":  homeDirExists,
		"isDir":        isDir,
	}
	return funcMap
}

func formatString(variable any) any {
	switch variable.(type) {
	case string, Template, Option:
		return fmt.Sprintf(`"%s"`, escapeString(variable))
	default:
		return variable
	}
}

func splitString(variable any) any {
	switch variable := variable.(type) {
	case string:
		variable = strings.TrimSpace(variable)
		if len(variable) == 0 {
			return []string{variable}
		}

		if strings.Contains(variable, "\n") {
			return strings.Split(variable, "\n")
		}

		return strings.Fields(variable)
	case Template:
		return splitString(variable.String())
	default:
		return variable
	}
}

func formatArray(variable any, delim ...string) any {
	delimiter := " "
	if len(delim) > 0 {
		delimiter = delim[0]
	}

	switch variable := variable.(type) {
	case string:
		split := splitString(variable).([]string)
		array := []string{}

		for _, value := range split {
			array = append(array, formatString(value).(string))
		}

		return strings.Join(array, delimiter)
	case Template:
		return formatArray(variable.String())
	default:
		return variable
	}
}

func escapeString(variable any) any {
	clean := func(v string) string {
		switch context.Current.Shell {
		case PWSH, POWERSHELL:
			return strings.NewReplacer(
				"`", "``",
				`"`, "`\"",
			).Replace(v)
		default:
			return strings.NewReplacer(
				`\`, `\\`,
				`"`, `\"`,
			).Replace(v)
		}
	}

	switch v := variable.(type) {
	case Template:
		value := v.String()
		return clean(value)
	case string:
		return clean(v)
	default:
		return variable
	}
}

func match(variable string, values ...string) bool {
	return slices.Contains(values, variable)
}

func hasCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// homeFileExists/homeDirExists intentionally resolve relative paths from .Home.
// Most template checks target files under the home directory, so this avoids requiring
// repetitive printf/path-join template expressions in common configurations.
func homeFileExists(path string) bool {
	info := pathExists(resolveFromHome(path))
	return info.exists && !info.isDir
}

func homeDirExists(path string) bool {
	info := pathExists(resolveFromHome(path))
	return info.exists && info.isDir
}

func pathExists(path string) pathInfo {
	pathExistsCacheMu.RLock()
	cached, OK := pathExistsCache[path]
	pathExistsCacheMu.RUnlock()
	if OK {
		return cached
	}

	info, err := os.Stat(path)
	result := pathInfo{
		exists: err == nil,
		isDir:  err == nil && info.IsDir(),
	}

	pathExistsCacheMu.Lock()
	pathExistsCache[path] = result
	pathExistsCacheMu.Unlock()

	return result
}

func resolveFromHome(path string) string {
	if filepath.IsAbs(path) || strings.HasPrefix(path, "/") {
		return path
	}

	return filepath.Join(context.Home(), path)
}

func clearPathExistsCache() {
	pathExistsCacheMu.Lock()
	defer pathExistsCacheMu.Unlock()
	pathExistsCache = map[string]pathInfo{}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
