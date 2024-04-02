package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type Aliae struct {
	Aliae   shell.Aliae   `yaml:"alias"`
	Envs    shell.Envs    `yaml:"env"`
	Paths   shell.Paths   `yaml:"path"`
	Scripts shell.Scripts `yaml:"script"`
}

func customUnmarshaler(a *Aliae, b []byte) error {
	data, err := includeUnmarshaler(b)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(data))
	if err = decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

// includeUnmarshaler handles unmarshaling of !include and !include_dir tags
func includeUnmarshaler(b []byte) ([]byte, error) {
	s := strings.Split(string(b), "\n")

	includeFuncMap := map[string]func([]string) (string, error){
		"!include_dir": getDirFiles,
		"!include":     getFile,
	}

	for i, line := range s {
		for key, f := range includeFuncMap {
			if !strings.HasPrefix(line, key) {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid %s directive: \n%s", key, line)
			}

			data, err := f(parts[1:])
			if err != nil {
				return nil, err
			}

			s[i] = data
			break
		}
	}

	returnData := []byte(strings.Join(s, "\n"))
	if len(returnData) == len(b) {
		return returnData, nil
	}

	return includeUnmarshaler(returnData)
}

func trimQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}

func getFile(s []string) (string, error) {
	// Allows for templating in the file path
	file := shell.Template(trimQuotes(strings.Join(s, ""))).Parse().String()

	// check if filepath is relative
	file, err := relativePath(file)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getDirFiles(d []string) (string, error) {
	// Allows for templating in the directory path
	dir := shell.Template(trimQuotes(strings.Join(d, ""))).Parse().String()

	// check if filepath is relative
	dir, err := relativePath(dir)
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var configData strings.Builder

	for _, file := range files {
		if !isYAMLExtension(file.Name()) {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}

		configData.WriteString("\n")
		configData.Write(data)
	}

	return configData.String(), nil
}

func relativePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	if len(configPathCache) == 0 {
		return "", errors.New("Config file not found")
	}

	if strings.HasPrefix(configPathCache, "https://") || strings.HasPrefix(configPathCache, "http://") {
		return "", errors.New("Remote files are not allowed to contain include directives")
	}

	// get the directory of the config file
	configPathCacheDir := filepath.Dir(configPathCache)

	// append the file to the directory
	path = filepath.Join(configPathCacheDir, path)

	return path, nil
}

func isYAMLExtension(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".yml" || ext == ".yaml"
}
