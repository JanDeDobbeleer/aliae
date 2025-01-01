package config

import (
	"bytes"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func templateUmarshaler(t *shell.Template, b []byte) error {
	return stringUmarshaler((*string)(t), b)
}

func stringUmarshaler(s *string, b []byte) error {
	if value, OK := unmarshalFoldedBlockScalar(string(b)); OK {
		*s = value
		return nil
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(b))
	return decoder.Decode(s)
}

func unmarshalFoldedBlockScalar(value string) (string, bool) {
	if !strings.HasPrefix(value, ">") {
		return value, false
	}

	lineBreak := "\n"
	if strings.Contains(value, "\r\n") {
		lineBreak = "\r\n"
	}

	value = strings.TrimLeft(value, ">")
	value = strings.TrimSpace(value)
	splitted := strings.Split(value, lineBreak)

	for i, line := range splitted {
		splitted[i] = strings.TrimSpace(line)
	}

	value = strings.Join(splitted, " ")

	return value, true
}
