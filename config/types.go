package config

// ClashConfig holds the fields extracted from the remote subscription.
// We use `any` to preserve the original YAML structure without losing fields.
type ClashConfig struct {
	Proxies     []any    `yaml:"proxies"`
	ProxyGroups []any    `yaml:"proxy-groups"`
	Rules       []string `yaml:"rules"`
}
