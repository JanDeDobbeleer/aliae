package shell

import (
	context_ "context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	NU              = "nu"
	NuEnvBlockStart = "export-env {\n"
	NuEnvBlockEnd   = "\n}"

	nuScriptName = "aliae.nu"

	// nuGuard makes the autoload script a no-op once aliae is uninstalled,
	// as removing the `aliae init nu` line doesn't remove the script.
	nuGuard = "if (which aliae | is-empty) { return }\n\n"
)

func (a *Alias) nu() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} = {{ .Value }}`
	case Function:
		a.template = `def {{ .Name }} [] {
    {{ .Value }}
}`
	}

	return a
}

func (e *Echo) nu() *Echo {
	e.template = `echo "{{ .Message }}"`
	return e
}

func (e *Env) nu() *Env {
	switch e.Type {
	case Array:
		e.template = `    $env.{{ .Name }} = [{{ formatArray .Value }}]`
	case String:
		fallthrough
	default:
		e.template = `    $env.{{ .Name }} = {{ formatString .Value }}`
	}

	return e
}

func (l *Link) nu() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }} out+err>| ignore`
	if context.Current.OS == context.WINDOWS {
		template = `{{ $source := (escapeString .Name) }}mklink {{ if isDir $source }}/d{{ else }}/h{{ end }} {{ $source }} {{ escapeString .Target }} out+err>| ignore`
	}

	l.template = template
	return l
}

func (p *Path) nu() *Path {
	template := `$env.%s = ($env.%s | prepend {{ formatString .Value }})`
	pathName := "PATH"

	if context.Current.OS == context.WINDOWS {
		pathName = "Path"
	}

	p.template = fmt.Sprintf(template, pathName, pathName)
	return p
}

// NuInit writes the init script to Nushell's vendor autoload directory.
// Nushell (v0.104.0+) loads autoload scripts after config.nu, so the script
// written by a single `aliae init nu` line in config.nu is picked up in the
// same session.
func NuInit(script string) error {
	removeLegacyNuScript()

	dir, err := nuAutoloadDir()
	if err != nil {
		return err
	}

	return writeFile(filepath.Join(dir, nuScriptName), []byte(nuGuard+script), 0o644)
}

// nuAutoloadDir resolves Nushell's vendor autoload directory. The result is
// cached on disk to avoid spawning nu on every shell startup.
func nuAutoloadDir() (string, error) {
	cachePath := nuAutoloadDirCachePath()

	if dir := cachedNuAutoloadDir(cachePath); len(dir) != 0 {
		return dir, nil
	}

	out, err := exec.CommandContext(context_.Background(), "nu", "-c", "$nu.data-dir | path join vendor autoload").Output()
	if err != nil {
		return "", fmt.Errorf("unable to resolve the Nushell autoload directory: %w", err)
	}

	dir := strings.TrimSpace(string(out))
	if len(dir) == 0 {
		return "", errors.New("unable to resolve the Nushell autoload directory")
	}

	if err = os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	cacheNuAutoloadDir(cachePath, dir)

	return dir, nil
}

func nuAutoloadDirCachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "aliae", "nu-autoload-dir")
}

func cachedNuAutoloadDir(cachePath string) string {
	if len(cachePath) == 0 {
		return ""
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return ""
	}

	dir := strings.TrimSpace(string(data))
	if len(dir) == 0 {
		return ""
	}

	// the autoload directory might have been removed since it was cached
	if _, err = os.Stat(dir); err != nil {
		return ""
	}

	return dir
}

func cacheNuAutoloadDir(cachePath, dir string) {
	if len(cachePath) == 0 {
		return
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o700); err != nil {
		return
	}

	_ = os.WriteFile(cachePath, []byte(dir), 0o644)
}

// removeLegacyNuScript deletes the pre-autoload init script at ~/.aliae.nu.
func removeLegacyNuScript() {
	_ = os.Remove(filepath.Join(context.Home(), ".aliae.nu"))
}
