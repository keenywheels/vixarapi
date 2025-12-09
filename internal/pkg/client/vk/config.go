package vk

import "time"

// HTTPConfig config for http client
type HTTPConfig struct {
	Timeout time.Duration `mapstructure:"timeout"`
}

// AuthConfig config for auth
type AuthConfig struct {
	BaseURL  string `mapstructure:"url"`
	ClientID string `mapstructure:"client_id"`
}

// Config VK client configuration
type Config struct {
	HTTP   HTTPConfig `mapstructure:"http"`
	Auth   AuthConfig `mapstructure:"auth"`
	Secret string     `mapstructure:"secret"`
}
