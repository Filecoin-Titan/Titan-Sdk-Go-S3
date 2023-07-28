package config

import (
	"time"
)

type Config struct {
	TitanAddress string
	Timeout      time.Duration
	CandidateID  string
}

// Option is a single titan sdk Config.
type Option func(opts *Config)

// DefaultOption returns a default set of options.
func DefaultOption() Config {
	return Config{
		TitanAddress: "",
		Timeout:      30 * time.Second,
	}
}

// TitanAddressOption set titan server address
func TitanAddressOption(address string) Option {
	return func(opts *Config) {
		opts.TitanAddress = address
	}
}

// TimeoutOption specifies a time limit for requests made by the http Client.
func TimeoutOption(timeout time.Duration) Option {
	return func(opts *Config) {
		opts.Timeout = timeout
	}
}

// CandidateIDOption specifies a candidate id
func CandidateIDOption(id string) Option {
	return func(opts *Config) {
		opts.CandidateID = id
	}
}
