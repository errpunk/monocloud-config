package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
	// Create a temp local config file (mimics example/config.yaml)
	localConfig := `tun:
  enable: true
  stack: system
dns:
  enable: true
allow-lan: true
external-controller: 0.0.0.0:9090
`
	tmpFile, err := os.CreateTemp("", "local-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(localConfig); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Build a fake remote config
	remote := &ClashConfig{
		Proxies: []any{
			map[string]any{"name": "HK1", "type": "ss", "server": "hk1.example.com", "port": 443},
		},
		ProxyGroups: []any{
			map[string]any{"name": "Proxy", "type": "select", "proxies": []any{"HK1"}},
		},
		Rules: []string{
			"DOMAIN-KEYWORD,google,Proxy",
			"MATCH,DIRECT",
		},
	}

	if err := Merge(tmpFile.Name(), remote); err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	// Read back and verify
	resultBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read merged file: %v", err)
	}

	var result map[string]any
	if err := yaml.Unmarshal(resultBytes, &result); err != nil {
		t.Fatalf("failed to parse merged YAML: %v", err)
	}

	// Check original fields preserved
	if result["allow-lan"] == nil {
		t.Error("expected 'allow-lan' to be preserved after merge")
	}
	if result["external-controller"] == nil {
		t.Error("expected 'external-controller' to be preserved after merge")
	}

	// Check remote fields merged
	proxies, ok := result["proxies"].([]any)
	if !ok || len(proxies) != 1 {
		t.Errorf("expected 1 proxy, got %v", result["proxies"])
	}

	rules, ok := result["rules"].([]any)
	if !ok || len(rules) != 2 {
		t.Errorf("expected 2 rules, got %v", result["rules"])
	}

	proxyGroups, ok := result["proxy-groups"].([]any)
	if !ok || len(proxyGroups) != 1 {
		t.Errorf("expected 1 proxy-group, got %v", result["proxy-groups"])
	}
}

func TestMerge_MissingLocalFile(t *testing.T) {
	remote := &ClashConfig{
		Proxies: []any{},
		Rules:   []string{},
	}
	err := Merge("/nonexistent/path/config.yaml", remote)
	if err == nil {
		t.Error("expected error for missing local config file")
	}
}
