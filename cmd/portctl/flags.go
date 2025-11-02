package main

import (
	"time"
)

// Config holds CLI options
type Config struct {
	Port       string
	Signal     string
	DryRun     bool
	Force      bool
	Timeout    time.Duration
	KillAfter  bool
	Verbose    bool
	Format     string // text | json | csv
}

// defaultConfig returns defaults
func defaultConfig() Config {
	return Config{
		Signal:    "TERM",
		DryRun:    false,
		Force:     false,
		Timeout:   5 * time.Second,
		KillAfter: false,
		Verbose:   false,
		Format:    "text",
	}
}
