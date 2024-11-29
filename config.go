package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Names have to be capitalized
type config struct {
	BaseURL      string `yaml:"baseURL"`
	Password     string `yaml:"password"`
	StatsdIPPort string `yaml:"statsdIPPort"`
}

// Embed the default configuration file
// Do not remove next line
//
//go:embed default_config.yml
var defaultConfig []byte

func loadConfig(configFile string) (*config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFile)

	// This applies to Linux/MacOS only, not Windows.
	permissions := os.FileMode(0600)

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File does not exist, write the default configuration
		err = os.WriteFile(configPath, defaultConfig, permissions) // Use the embedded defaultConfig
		if err != nil {
			return nil, fmt.Errorf("failed to write default config file: %w", err)
		}
		log.Printf("Default config written to %s", configPath)
	}

	// Open the config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Decode the config file
	config := &config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// determineConfigFile sets configFile based on GOOS
func determineConfigFile() string {
	configFilename := "att-fiber-gateway-info.yml"
	var configFile string

	switch runtime.GOOS {
	case "windows":
		configFile = configFilename
	default: // Linux and other Unix-like systems
		configFile = "." + configFilename
	}

	return configFile
}

// loadAppConfig loads configuration from the file
func loadAppConfig(configFile string) *config {
	config, err := loadConfig(configFile)
	if err != nil {
		logFatalf("Error loading configuration: %v", err)
	}
	return config
}
