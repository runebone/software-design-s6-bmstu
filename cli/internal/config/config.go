package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log        LogConfig        `toml:"log"`
	Client     ClientConfig     `toml:"cli"`
	Aggregator AggregatorConfig `toml:"aggregator"`
}

type LogConfig struct {
	Path  string `toml:"path"`
	Level string `toml:"level"`
}

type ClientConfig struct {
	Path          string    `toml:"path"`
	ContainerName string    `toml:"container_name"`
	TokensPath    string    `toml:"tokens_path"`
	Log           LogConfig `toml:"log"`
}

type AggregatorConfig struct {
	Path          string    `toml:"path"`
	ContainerName string    `toml:"container_name"`
	BaseURL       string    `toml:"base_url"`
	LocalPort     int       `toml:"local_port"`
	ExposedPort   int       `toml:"exposed_port"`
	Log           LogConfig `toml:"log"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}

	s := &config.Client.Log
	s.Path = config.Log.Path + "/" + s.Path
	if s.Level == "" {
		s.Level = config.Log.Level
	}

	return &config, nil
}
