package config

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/jandedobbeleer/aliae/src/shell"
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
)

type Aliae struct {
	Aliae []*shell.Alias `yaml:"aliae"`
}

func LoadConfig(configPath string) (*Aliae, error) {
	if strings.HasPrefix(configPath, "https://") {
		return getRemoteConfig(configPath)
	}

	if len(configPath) == 0 {
		configPath = path.Join(home(), ".aliae.yaml")
	}

	if filepath, err := os.Stat(configPath); os.IsNotExist(err) || filepath.IsDir() {
		return nil, fmt.Errorf("Config file not found: %s", configPath)
	}

	data, _ := os.ReadFile(configPath)

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

func getRemoteConfig(url string) (*Aliae, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to download config file: %s\n→ %s", url, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

func parseConfig(data []byte) (*Aliae, error) {
	var aliae Aliae
	err := yaml.Unmarshal(data, &aliae)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %s", err)
	}
	return &aliae, nil
}
