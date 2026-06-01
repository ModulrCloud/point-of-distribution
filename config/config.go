package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WSInterface        string `json:"wsInterface"`
	WSPort             int    `json:"wsPort"`
	DataPath           string `json:"dataPath"`
	MaxConcurrentLocks int    `json:"maxConcurrentLocks"`
	LogRequests        bool   `json:"logRequests"`
}

func Load() Config {
	cfg := Config{
		WSInterface:        "0.0.0.0",
		WSPort:             9070,
		DataPath:           "poddata",
		MaxConcurrentLocks: 100,
		LogRequests:        true,
	}

	configPath := os.Getenv("POD_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs.json"
	}
	if !filepath.IsAbs(configPath) {
		cwd, _ := os.Getwd()
		configPath = filepath.Join(cwd, configPath)
	}

	if data, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(data, &cfg)
	}

	if !filepath.IsAbs(cfg.DataPath) {
		cwd, _ := os.Getwd()
		cfg.DataPath = filepath.Join(cwd, cfg.DataPath)
	}

	_ = os.MkdirAll(cfg.DataPath, os.ModePerm)

	return cfg
}
