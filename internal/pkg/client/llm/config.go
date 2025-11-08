package llm

import "time"

// Config LLM client configuration
type Config struct {
	Token   string        `mapstructure:"token"`
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
}
