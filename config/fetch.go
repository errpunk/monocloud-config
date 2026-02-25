package config

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// Fetch downloads the Clash subscription from MONOCLOUD_URL env var and returns
// the extracted proxies, proxy-groups, and rules.
func Fetch() (*ClashConfig, error) {
	url := os.Getenv("MONOCLOUD_URL")
	if url == "" {
		return nil, fmt.Errorf("environment variable MONOCLOUD_URL is not set")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	// monocloud requires clash.meta as user-agent, otherwise it returns an error
	req.Header.Set("User-Agent", "clash.meta")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download subscription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subscription download failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var cfg ClashConfig
	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse remote YAML: %w", err)
	}

	if len(cfg.Proxies) == 0 && len(cfg.ProxyGroups) == 0 && len(cfg.Rules) == 0 {
		return nil, fmt.Errorf("remote config appears empty or invalid (no proxies, proxy-groups, or rules found)")
	}

	return &cfg, nil
}
