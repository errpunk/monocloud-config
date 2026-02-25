package config

// ClashConfig represents the structure of a Clash/Mihomo configuration file.
// We use `any` for proxies and proxy-groups to preserve the original YAML structure
// without losing any fields.
type ClashConfig struct {
	Proxies     []any    `yaml:"proxies,omitempty"`
	ProxyGroups []any    `yaml:"proxy-groups,omitempty"`
	Rules       []string `yaml:"rules,omitempty"`
}
