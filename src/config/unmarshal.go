package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type Aliae struct {
	Aliae   shell.Aliae   `yaml:"alias"`
	Envs    shell.Envs    `yaml:"env"`
	Paths   shell.Paths   `yaml:"path"`
	CDPaths shell.CDPaths `yaml:"cdpath"`
	Scripts shell.Scripts `yaml:"script"`
	Links   shell.Links   `yaml:"link"`
}

// sectionKeys are the top-level Aliae document keys (config.go's Aliae struct tags)
// that an included file may itself be wrapped in, e.g. a file combining alias/env/script
// sections that gets included separately at each section's own position.
var sectionKeys = map[string]bool{
	"alias":  true,
	"env":    true,
	"path":   true,
	"cdpath": true,
	"script": true,
	"link":   true,
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

	tagNames := [][]byte{
		[]byte("!include_dir"),
		[]byte("!include"),
	}

	for i, line := range s {
		for _, name := range tagNames {
			if !bytes.Contains(line, name) {
				continue
			}

			parts := bytes.Fields(line)
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid %s directive: \n%s", name, line)
			}

			tagIdx := bytes.Index(line, name)
			argument := bytes.TrimSpace(line[tagIdx+len(name):])

			content, ok := unquoteArgument(argument)
			if !ok {
				return nil, fmt.Errorf("invalid %s directive: \n%s", name, line)
			}

			folder, condition, hasCondition := splitIncludeCondition(content)

			if hasCondition && shell.If(condition).Ignore() {
				s[i] = skippedIncludeLine(parts[0])
				break
			}

			path, err := validatePath(folder)
			if err != nil {
				return nil, err
			}

			destKey := sectionDestKey(parts[0])
			isDir := string(name) == "!include_dir"

			var data []byte
			if isDir {
				data, err = readDir(path, destKey)
			} else {
				data, err = os.ReadFile(path)
			}

			if err != nil {
				return nil, err
			}

			if !isDir && destKey != "" {
				extracted, wrapped, extractErr := extractSection(data, destKey)
				if extractErr != nil {
					return nil, extractErr
				}

				if wrapped && len(extracted) == 0 {
					s[i] = skippedIncludeLine(parts[0])
					break
				}

				if wrapped {
					data = extracted
				}
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

// unquoteArgument strips one matching pair of outer quotes from an include directive's
// argument, e.g. !include "path.yaml" or !include 'path.yaml if="hasCommand kubectl"'.
// The condition (if any) must live inside these same outer quotes: the raw YAML document
// must be syntactically valid before any of this preprocessing ever runs, so nothing may
// follow the quoted argument on the line. A double-quoted argument supports `\"` to embed a
// literal double quote; a single-quoted argument treats its content literally.
func unquoteArgument(argument []byte) (content string, ok bool) {
	if len(argument) == 0 {
		return "", false
	}

	quote := argument[0]
	if quote != '"' && quote != '\'' {
		// bare, unquoted argument: a plain YAML scalar can't carry a trailing
		// condition (that would be a second value on the same line), so the
		// whole argument is the path.
		return string(argument), true
	}

	if len(argument) < 2 {
		return "", false
	}

	if quote == '\'' {
		if argument[len(argument)-1] != '\'' {
			return "", false
		}

		return string(argument[1 : len(argument)-1]), true
	}

	var builder strings.Builder

	for idx := 1; idx < len(argument)-1; idx++ {
		if argument[idx] == '\\' && idx+1 < len(argument)-1 && argument[idx+1] == '"' {
			builder.WriteByte('"')
			idx++

			continue
		}

		builder.WriteByte(argument[idx])
	}

	if argument[len(argument)-1] != '"' {
		return "", false
	}

	return builder.String(), true
}

// splitIncludeCondition splits an unquoted include argument into its path and an optional
// trailing `if=<condition>` marker, e.g. `path.yaml if=hasCommand "kubectl"`.
func splitIncludeCondition(content string) (path, condition string, ok bool) {
	marker := " if="

	before, after, ok := strings.Cut(content, marker)
	if !ok {
		return content, "", false
	}

	return before, after, true
}

// skippedIncludeLine replaces a skipped (condition evaluated to false) include directive
// with a YAML-safe no-op: `null` for a map key, or nothing for a list item.
func skippedIncludeLine(key []byte) []byte {
	if string(key) == "-" {
		return []byte{}
	}

	result := key[0:]
	result = append(result, []byte(" null")...)

	return result
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

// readDir reads every YAML file in dir and joins their content with a blank line.
// When destKey is one of sectionKeys, each file's content is first passed through
// extractSection: a file wrapped in that section contributes only that section
// (or nothing, if it doesn't carry that section); a bare-list file (today's shape)
// is used unchanged.
func readDir(dir, destKey string) ([]byte, error) {
	files, err := os.ReadDir(dir)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		// If the directory does not exist, treat it as empty
		return []byte{}, nil
	case err != nil:
		return []byte{}, err
	}

	var chunks [][]byte

	for _, file := range files {
		if !isYAMLExtension(file.Name()) {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		if destKey != "" {
			extracted, wrapped, err := extractSection(data, destKey)
			if err != nil {
				return nil, err
			}

			if wrapped {
				data = extracted
			}
		}

		if len(data) == 0 {
			continue
		}

		chunks = append(chunks, data)
	}

	return bytes.Join(chunks, []byte("\n")), nil
}

// sectionDestKey returns the top-level Aliae section (see sectionKeys) an include
// directive targets directly, e.g. "alias" for `alias: !include ...`. It returns ""
// for a list-item include (`- !include ...`, no ancestor key tracked) or a key that
// isn't a recognized section — signaling that section extraction must not run.
func sectionDestKey(key []byte) string {
	trimmed := strings.TrimSuffix(string(key), ":")
	if !sectionKeys[trimmed] {
		return ""
	}

	return trimmed
}

// extractSection inspects an included file's content for a top-level mapping
// carrying one or more sectionKeys, e.g. a file combining alias/env/script into one
// document. When such a mapping is found, wrapped is true and extracted holds the
// raw source text of destKey's value (empty when the file doesn't carry that
// section), so nested !include tags, block scalars, and Go-template syntax within
// it survive untouched for the next unmarshal pass. wrapped is false when the
// content isn't shaped this way (e.g. a bare list), signaling the caller to use the
// original content unchanged.
func extractSection(data []byte, destKey string) (extracted []byte, wrapped bool, err error) {
	file, err := parser.ParseBytes(data, 0)
	if err != nil {
		// not valid enough to inspect; let the normal decode path surface the error
		return nil, false, nil
	}

	if len(file.Docs) == 0 || file.Docs[0].Body == nil {
		return nil, false, nil
	}

	if len(file.Docs) > 1 {
		return nil, false, errors.New("included file with multiple YAML documents is not supported")
	}

	mapping, ok := file.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return nil, false, nil
	}

	var found *ast.MappingValueNode

	for _, value := range mapping.Values {
		key := strings.Trim(value.Key.String(), `"'`)
		if !sectionKeys[key] {
			continue
		}

		wrapped = true

		if key == destKey {
			found = value
		}
	}

	if !wrapped || found == nil {
		return nil, wrapped, nil
	}

	return []byte(found.Value.String()), true, nil
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
