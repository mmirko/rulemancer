/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	debugLevelMax = 20
)

type Config struct {
	ClipsLessMode bool     `json:"clipsless_mode"`
	Debug         bool     `json:"debug"`
	DebugLevel    int      `json:"debug_level"`
	TLSCertFile   string   `json:"tls_cert_file"`
	TLSKeyFile    string   `json:"tls_key_file"`
	Games         []string `json:"games"`
}

func NewConfig() *Config {
	return &Config{
		ClipsLessMode: false,
		Debug:         false,
		TLSCertFile:   "server.crt",
		TLSKeyFile:    "server.key",
		Games:         []string{},
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
	if c.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/SaveConfig]")+" ", 0)
		l.Printf("Saved configuration to %s", path)
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
	if c.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/LoadConfig]")+" ", 0)
		l.Printf("Loading configuration from %s", path)
	}
	return nil
}
