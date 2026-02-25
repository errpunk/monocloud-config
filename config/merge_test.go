package config

import (
	"os"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestMerge(t *testing.T) {
	// Local config with specific key order we want to preserve
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

	resultBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read merged file: %v", err)
	}

	// Parse with MapSlice to verify order
	var result yaml.MapSlice
	if err := yaml.Unmarshal(resultBytes, &result); err != nil {
		t.Fatalf("failed to parse merged YAML: %v", err)
	}

	// Collect keys in order
	keys := make([]string, 0, len(result))
	for _, item := range result {
		keys = append(keys, item.Key.(string))
	}

	// Original keys must come first
	expectedPrefix := []string{"tun", "dns", "allow-lan", "external-controller"}
	for i, k := range expectedPrefix {
		if i >= len(keys) || keys[i] != k {
			t.Errorf("key order wrong at position %d: want %q, got %v", i, k, keys)
		}
	}

	// Remote keys must be at the end
	expectedSuffix := []string{"proxies", "proxy-groups", "rules"}
	offset := len(keys) - len(expectedSuffix)
	for i, k := range expectedSuffix {
		if keys[offset+i] != k {
			t.Errorf("remote key order wrong at position %d: want %q, got %q", offset+i, k, keys[offset+i])
		}
	}
}

func TestMerge_MissingLocalFile(t *testing.T) {
	remote := &ClashConfig{Proxies: []any{}, Rules: []string{}}
	err := Merge("/nonexistent/path/config.yaml", remote)
	if err == nil {
		t.Error("expected error for missing local config file")
	}
}
