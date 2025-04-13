package config

import (
	"bytes"
	context_ "context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/jandedobbeleer/aliae/src/context"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	defaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	client httpClient = &http.Client{Transport: defaultTransport}

	configPathCache string
)

func LoadConfig(configPath string) (*Aliae, error) {
	configPathCache = resolveConfigPath(configPath)

	if strings.HasPrefix(configPathCache, "http://") || strings.HasPrefix(configPathCache, "https://") {
		return getRemoteConfig(configPathCache)
	}

	if filepath, err := os.Stat(configPathCache); os.IsNotExist(err) || filepath.IsDir() {
		return nil, fmt.Errorf("config file not found: %s", configPathCache)
	}

	data, _ := os.ReadFile(configPathCache)

	return parseConfig(data)
}

func home() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}

	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}

	return home
}

func resolveConfigPath(configPath string) string {
	if len(configPath) == 0 {
		configPath = os.Getenv("ALIAE_CONFIG")
	}

	if len(configPath) == 0 {
		configPath = path.Join(home(), ".aliae.yaml")
	}

	return replaceTildePrefixWithHomeDir(configPath)
}

func replaceTildePrefixWithHomeDir(dir string) string {
	if !strings.HasPrefix(dir, "~") {
		return dir
	}

	rem := dir[1:]
	if len(rem) == 0 || isSeparator(rem[0]) {
		return home() + rem
	}

	return dir
}

func isSeparator(c uint8) bool {
	if c == '/' {
		return true
	}

	if runtime.GOOS == context.WINDOWS && c == '\\' {
		return true
	}

	return false
}

func getRemoteConfig(url string) (*Aliae, error) {
	req, err := http.NewRequestWithContext(context_.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download config file: %s\n→ %s", url, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

func parseConfig(data []byte) (*Aliae, error) {
	var aliae Aliae

	decoder := yaml.NewDecoder(bytes.NewBuffer(data), yaml.CustomUnmarshaler(aliaeUnmarshaler))
	err := decoder.Decode(&aliae)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return &aliae, nil
}
