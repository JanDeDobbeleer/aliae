package shell

import (
	"bytes"
	context_ "context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"text/template"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Template string

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
		"isPwshOption":   isPwshOption,
		"isPwshScope":    isPwshScope,
		"formatString":   formatString,
		"formatArray":    formatArray,
		"escapeString":   escapeString,
		"env":            os.Getenv,
		"match":          match,
		"hasCommand":     hasCommand,
		"isDir":          isDir,
		"dirAccessible":  dirAccessible,
		"pathAccessible": pathAccessible,
		"wslPath":        wslPath,
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

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// pathAccessible reports whether path exists and is readable by the current user.
func pathAccessible(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return canReadPath(path)
}

// dirAccessible reports whether path exists, is a directory, and is traversable
// (i.e. its contents can be listed) by the current user.
func dirAccessible(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return canTraverseDir(path)
}

// wslPath converts a Windows path to its WSL equivalent via the wslpath binary.
// wslpath only exists inside WSL's interop layer, so its presence on PATH is
// itself the WSL signal; path is returned unchanged when the binary isn't found
// or fails to convert it (e.g. running natively, or the path can't be resolved).
func wslPath(path string) string {
	bin, err := exec.LookPath("wslpath")
	if err != nil {
		return path
	}

	out, err := exec.CommandContext(context_.Background(), bin, "-a", path).Output()
	if err != nil {
		return path
	}

	return strings.TrimSpace(string(out))
}
