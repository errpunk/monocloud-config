package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// remoteKeys are the keys that will be stripped from the local config
// and re-appended at the end with values from the remote subscription.
var remoteKeys = []string{"proxies", "proxy-groups", "rules"}

// Merge reads the local config file at configPath, removes any existing
// proxies/proxy-groups/rules entries (preserving all other keys and their
// original order), then appends the remote values at the end.
func Merge(configPath string, remote *ClashConfig) error {
	localBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read local config %q: %w", configPath, err)
	}

	// MapSlice preserves insertion order — critical for stable YAML output.
	var localMap yaml.MapSlice
	if err := yaml.Unmarshal(localBytes, &localMap); err != nil {
		return fmt.Errorf("failed to parse local config YAML: %w", err)
	}

	// Remove any existing proxies / proxy-groups / rules so we can
	// re-append them at the very end.
	skipSet := make(map[string]bool, len(remoteKeys))
	for _, k := range remoteKeys {
		skipSet[k] = true
	}
	filtered := localMap[:0]
	for _, item := range localMap {
		if key, ok := item.Key.(string); ok && skipSet[key] {
			continue
		}
		filtered = append(filtered, item)
	}

	// Append remote fields at the end, in a fixed order.
	filtered = append(filtered,
		yaml.MapItem{Key: "proxies", Value: remote.Proxies},
		yaml.MapItem{Key: "proxy-groups", Value: remote.ProxyGroups},
		yaml.MapItem{Key: "rules", Value: remote.Rules},
	)

	out, err := yaml.Marshal(filtered)
	if err != nil {
		return fmt.Errorf("failed to marshal merged config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0o644); err != nil {
		return fmt.Errorf("failed to write merged config to %q: %w", configPath, err)
	}

	return nil
}
