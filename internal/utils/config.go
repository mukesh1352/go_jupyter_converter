package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DatasetRoot string `json:"dataset_root"`
	OutputDir   string `json:"output_dir"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "tool")
	configPath := filepath.Join(configDir, "config.json")

	// Auto-create if missing
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file not found. Creating default config at:", configPath)

		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create config directory: %v", err)
		}

		defaultCfg := &Config{
			DatasetRoot: filepath.Join(homeDir, "datasets"),
			OutputDir:   filepath.Join(homeDir, "output"),
		}

		data, err := json.MarshalIndent(defaultCfg, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("could not serialize default config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("could not write config file: %v", err)
		}

		return defaultCfg, nil
	}

	// Read config if it exists
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config format: %v", err)
	}

	return &cfg, nil
}
