package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type VideoMode string

// Update the supported modes in main.go upon any addition/removal to the constants below.
const (
	VideoModeMovie  VideoMode = "movie"
	VideoModeSeries VideoMode = "series"
)

type Config struct {
	Mode              VideoMode `yaml:"mode"`
	EpisodesPerSeason int       `yaml:"episodes_per_season"`
}

// NewConfig returns the default configuration when the specified file does not exist.
func NewConfig(path string) (*Config, error) {
	bs, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file: %w", err)
	}

	var c Config
	if err := yaml.Unmarshal(bs, &c); err != nil {
		return nil, fmt.Errorf("failed to parse the config file: %w", err)
	}
	return &c, nil
}

func defaultConfig() *Config {
	return &Config{
		Mode: VideoModeMovie,
	}
}
