package config

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

// aliDNSClient is an HTTP client that forces DNS resolution through
// Alibaba Public DNS (223.5.5.5), bypassing any system or VPN DNS.
var aliDNSClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, "udp", "223.5.5.5:53")
				},
			},
		}).DialContext,
	},
}

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

	resp, err := aliDNSClient.Do(req)
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
