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
	Links   shell.Links   `yaml:"link"`
}

type FuncMap []StringFunc
type StringFunc struct {
	F    func(string) ([]byte, error)
	Name []byte
}

func aliaeUnmarshaler(a *Aliae, b []byte) error {
	data, err := includeUnmarshaler(b)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(data), yaml.CustomUnmarshaler(templateUmarshaler))
	if err = decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

// includeUnmarshaler handles unmarshaling of !include and !include_dir tags
func includeUnmarshaler(b []byte) ([]byte, error) {
	newline := []byte("\n")

	s := bytes.Split(b, newline)

	includeFuncMap := FuncMap{
		{
			Name: []byte("!include_dir"),
			F:    readDir,
		},
		{
			Name: []byte("!include"),
			F:    os.ReadFile,
		},
	}

	for i, line := range s {
		for _, f := range includeFuncMap {
			if !bytes.Contains(line, f.Name) {
				continue
			}

			parts := bytes.Fields(line)
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid %s directive: \n%s", f.Name, line)
			}

			folder := string(bytes.Join(parts[2:][0:], []byte(" ")))
			path, err := validatePath(folder)
			if err != nil {
				return nil, err
			}

			data, err := f.F(path)
			if err != nil {
				return nil, err
			}

			splitted := bytes.Split(data, newline)
			for i, line := range splitted {
				splitted[i] = indent(line)
			}

			indented := bytes.Join(splitted, newline)

			result := parts[0][0:]

			switch string(result) {
			case "-":
				// check if we're in the list instead of the key
				// if so, drop the dash and start with a newline
				result = newline
			default:
				result = append(result, newline...)
			}

			result = append(result, indented...)

			s[i] = result
			break
		}
	}

	data := bytes.Join(s, newline)
	if len(data) == len(b) {
		return data, nil
	}

	return includeUnmarshaler(data)
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

func indent(data []byte) []byte {
	newData := make([]byte, len(data)+2)
	newData[0] = ' '
	newData[1] = ' '

	copy(newData[2:], data)

	return newData
}

func readDir(dir string) ([]byte, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return []byte{}, err
	}

	var configData []byte

	for i, file := range files {
		if !isYAMLExtension(file.Name()) {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		configData = append(configData, data...)

		if i != len(files)-1 {
			configData = append(configData, []byte("\n")...)
		}
	}

	return configData, nil
}

func validatePath(path string) (string, error) {
	// Allows for templating in the file path
	path = shell.Template(trimQuotes(path)).Parse().String()

	if filepath.IsAbs(path) {
		return path, nil
	}

	if len(configPathCache) == 0 {
		return "", errors.New("config file not found")
	}

	if strings.HasPrefix(configPathCache, "https://") || strings.HasPrefix(configPathCache, "http://") {
		return "", errors.New("remote files are not allowed to contain include directives")
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
