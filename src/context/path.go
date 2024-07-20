package context

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Path []string

func getPath() *Path {
	if Current != nil {
		return Current.Path
	}

	path := &Path{}
	paths := os.Getenv("PATH")

	for _, p := range strings.Split(paths, PathDelimiter()) {
		path.Append(cleanPath(p))
	}

	return path
}

func (p *Path) Append(path string) {
	if len(path) == 0 || p.Contains(cleanPath(path)) {
		return
	}

	current := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s%s%s", path, PathDelimiter(), current))

	*p = append(*p, path)
}

func (p *Path) Contains(path string) bool {
	return slices.Contains(*p, cleanPath(path))
}

func cleanPath(path string) string {
	return strings.TrimRight(path, PathSeparator())
}
