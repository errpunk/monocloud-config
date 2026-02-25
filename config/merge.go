package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Merge reads the local config file at configPath, replaces proxies/proxy-groups/rules
// with those from remote, and writes the result back to configPath.
func Merge(configPath string, remote *ClashConfig) error {
	// Read local config as a generic map to preserve all existing fields
	localBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read local config %q: %w", configPath, err)
	}

	var localMap map[string]any
	if err := yaml.Unmarshal(localBytes, &localMap); err != nil {
		return fmt.Errorf("failed to parse local config YAML: %w", err)
	}

	if localMap == nil {
		localMap = make(map[string]any)
	}

	// Overwrite with remote values
	localMap["proxies"] = remote.Proxies
	localMap["proxy-groups"] = remote.ProxyGroups
	localMap["rules"] = remote.Rules

	out, err := yaml.Marshal(localMap)
	if err != nil {
		return fmt.Errorf("failed to marshal merged config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0o644); err != nil {
		return fmt.Errorf("failed to write merged config to %q: %w", configPath, err)
	}

	return nil
}
