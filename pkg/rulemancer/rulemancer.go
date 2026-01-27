/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ClipsLessMode bool                `json:"clipsless_mode"`
	Debug         bool                `json:"debug"`
	TLSCertFile   string              `json:"tls_cert_file"`
	TLSKeyFile    string              `json:"tls_key_file"`
	RulePool      string              `json:"rule_pool"`
	Assertables   []string            `json:"assertables"`
	Results       map[string][]string `json:"results"`
	Querables     []string            `json:"querables"`
}

func NewConfig() *Config {
	return &Config{
		ClipsLessMode: false,
		Debug:         false,
		TLSCertFile:   "server.crt",
		TLSKeyFile:    "server.key",
		RulePool:      "rulepool",
	}
}

// Save the current configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	// Check if the path is already existing
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	configBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(path, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}
	return nil
}

// LoadConfig loads the configuration from a JSON file
func (c *Config) LoadConfig(path string) error {
	// Check if the path is not existing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// ShowConfig prints the current configuration to the console
func (c *Config) ShowConfig() {
	fmt.Println("Current configuration:")
	fmt.Printf("Debug: %v\n", c.Debug)
}

// responseForType returns the list of status queries for a given assert type
func (c *Config) responseForType(assertType string) ([]string, error) {
	if c == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if val, ok := c.Results[assertType]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("unknown assert type: %s", assertType)
	}
}
