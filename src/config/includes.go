package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jandedobbeleer/aliae/src/shell"
	"gopkg.in/yaml.v3"
)

type IncludeTagType int

const (
	IncludeTagFile IncludeTagType = iota
	IncludeTagDir
)

var includeTagProcessors = map[string]func(*yaml.Node, string) error{
	"!include":           func(n *yaml.Node, p string) error { return processIncludeTag(n, p, IncludeTagFile, readRawYAMLFile) },
	"!include_dir":       func(n *yaml.Node, p string) error { return processIncludeTag(n, p, IncludeTagDir, readRawYAMLFile) },
	"!include_dir_list":  func(n *yaml.Node, p string) error { return processIncludeTag(n, p, IncludeTagDir, formatAsIncludeList) },
	"!include_dir_named": func(n *yaml.Node, p string) error { return processIncludeTag(n, p, IncludeTagDir, formatAsIncludeMap) },
}

// Walk YAML node tree and resolve and load include files
func resolveIncludes(node *yaml.Node, defaultPath string) error {
	if node == nil {
		return nil
	}

	if processor, found := includeTagProcessors[node.Tag]; found {
		err := processor(node, defaultPath)
		if err != nil {
			return err
		}
		return resolveIncludes(node, defaultPath)
	}

	for _, child := range node.Content {
		if err := resolveIncludes(child, defaultPath); err != nil {
			return err
		}
	}

	return nil
}

func processIncludeTag(node *yaml.Node, defaultPath string, includeType IncludeTagType, tagFormatter IncludeTagFormatter) error {
	path, err := validateAbsolutePath(node.Value, defaultPath, includeType)
	if err != nil {
		return err
	}

	filePaths, err := collectIncludePaths(path, includeType)
	if err != nil {
		return err
	}

	var includeDatas [][]byte
	for _, path := range filePaths {
		data, err := tagFormatter(path)
		if err != nil {
			return err
		}
		includeDatas = append(includeDatas, data)
	}

	finalData := bytes.Join(includeDatas, []byte("\n"))

	var includeNode yaml.Node
	if err := yaml.Unmarshal(finalData, &includeNode); err != nil {
		return fmt.Errorf("failed to parse included YAML: %w", err)
	}

	if err := unwrapRootNode(&includeNode); err != nil {
		return err
	}

	*node = includeNode
	return nil
}

func collectIncludePaths(path string, includeType IncludeTagType) ([]string, error) {
	switch includeType {
	case IncludeTagFile:
		return []string{path}, nil
	case IncludeTagDir:
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
		}

		filePaths := make([]string, 0, len(entries))
		for _, entry := range entries {
			if !entry.IsDir() && isYAMLExtension(entry.Name()) {
				filePaths = append(filePaths, filepath.Join(path, entry.Name()))
			}
		}
		sort.Strings(filePaths)
		return filePaths, nil
	}
	return nil, fmt.Errorf("unexpected include type: %v", includeType)
}

type IncludeTagFormatter func(filePath string) ([]byte, error)

func formatAsIncludeList(filePath string) ([]byte, error) {
	return []byte(fmt.Sprintf("- !include %q", filePath)), nil
}

func formatAsIncludeMap(filePath string) ([]byte, error) {
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return []byte(fmt.Sprintf("%q: !include %q", baseName, filePath)), nil
}

func readRawYAMLFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	data = bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))

	return data, nil
}

func unwrapRootNode(rootNode *yaml.Node) error {
	if rootNode == nil {
		return fmt.Errorf("no YAML node provided to unwrap tree")
	}

	// When empty yaml files are loaded, the root node is an unknown Kind (0) with !!null tag
	// Substitute it for a null scalar node.
	if rootNode.Kind == 0 { // Empty file scenario
		fmt.Println("For Testing !!!!! Warning: Node is unknown, setting as null")
		*rootNode = yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!null",
			Value: "",
		}
		return nil
	}

	if rootNode.Kind == yaml.DocumentNode && len(rootNode.Content) == 1 {
		// Unwrap document to first child
		*rootNode = *rootNode.Content[0]
		return nil
	}

	return fmt.Errorf("unexpected condition in YAML document node")
}

func validateAbsolutePath(path string, defaultPath string, includeType IncludeTagType) (string, error) {
	path = shell.Template(path).Parse().String()

	if !filepath.IsAbs(path) {
		if defaultPath == "" {
			return "", errors.New("default path not provided")
		}
		path = filepath.Join(defaultPath, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("path does not exist or is inaccessible: %w", err)
	}

	switch includeType {
	case IncludeTagFile:
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("expected file but got directory: %s", path)
		}
		if !isYAMLExtension(path) {
			return "", fmt.Errorf("invalid file extension for %s: only .yaml and .yml are supported", path)
		}
	case IncludeTagDir:
		if !info.Mode().IsDir() {
			return "", fmt.Errorf("expected directory but got a path: %s", path)
		}
	}

	return path, nil
}

func isYAMLExtension(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".yml" || ext == ".yaml"
}

// *********************************************************************************
// TRACE
// *********************************************************************************
var yamlNodeKinds = map[yaml.Kind]string{
	0:                 "Unknown",
	yaml.DocumentNode: "DocumentNode",
	yaml.SequenceNode: "SequenceNode",
	yaml.MappingNode:  "MappingNode",
	yaml.ScalarNode:   "ScalarNode",
	yaml.AliasNode:    "AliasNode",
}

func printYAMLNode(node *yaml.Node, indent int) {
	kindStr := yamlNodeKinds[node.Kind]

	indentStr := strings.Repeat("  ", indent) // 2-space indentation

	// Avoid initial indentation on the first line
	if indent == 0 {
		fmt.Printf("- Kind: %s(%d), Tag: %s, Value: %s\n", kindStr, int(node.Kind), node.Tag, node.Value)
	} else {
		fmt.Printf("%s- Kind: %s(%d), Tag: %s, Value: %s\n", indentStr, kindStr, int(node.Kind), node.Tag, node.Value)
	}

	if node.Content != nil {
		for _, child := range node.Content {
			printYAMLNode(child, indent+1)
		}
	}
}