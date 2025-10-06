package app

import (
	"flag"
	"os"
)

// default opts values
const (
	defaultConfigPath = "./configs/config.yaml"
)

// envs
const (
	envConfigPath = "CONFIG_PATH"
)

// Options represents application's options
type Options struct {
	ConfigPath string
}

// NewDefaultOpts creates default options
func NewDefaultOpts() *Options {
	return &Options{
		ConfigPath: defaultConfigPath,
	}
}

// LoadEnv updates options with values from envs
func (opts *Options) LoadEnv() {
	if val, ok := os.LookupEnv(envConfigPath); ok {
		opts.ConfigPath = val
	}
}

// LoadFlags updates options with values from cmd flags
func (opts *Options) LoadFlags() {
	configFlag := flag.String("config", opts.ConfigPath, "path to config file")
	flag.Parse()

	opts.ConfigPath = *configFlag
}
