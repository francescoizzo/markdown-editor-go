package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration
type Config struct {
	// Theme settings
	IsDarkMode bool `json:"isDarkMode"`

	// Editor settings
	FontSize    int    `json:"fontSize"`
	FontFamily  string `json:"fontFamily"`
	TabSize     int    `json:"tabSize"`
	LineNumbers bool   `json:"lineNumbers"`

	// Autosave settings
	AutoSaveEnabled bool `json:"autoSaveEnabled"`
	AutoSaveDelay   int  `json:"autoSaveDelay"` // in seconds

	// Recent files
	RecentFiles []string `json:"recentFiles"`

	// Window settings
	WindowWidth  int `json:"windowWidth"`
	WindowHeight int `json:"windowHeight"`

	// File path of the configuration file itself
	configPath string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		IsDarkMode:      false,
		FontSize:        14,
		FontFamily:      "Roboto Mono, monospace",
		TabSize:         4,
		LineNumbers:     true,
		AutoSaveEnabled: true,
		AutoSaveDelay:   5, // 5 seconds
		RecentFiles:     []string{},
		WindowWidth:     1024,
		WindowHeight:    768,
	}
}

// Load loads the configuration from file
func (c *Config) Load() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	c.configPath = configPath

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := DefaultConfig()
		defaultConfig.configPath = configPath
		return defaultConfig.Save()
	}

	// Read config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Parse JSON
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	if c.configPath == "" {
		configPath, err := getConfigPath()
		if err != nil {
			return err
		}
		c.configPath = configPath
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(c.configPath, data, 0644)
}

// AddRecentFile adds a file to the recent files list
func (c *Config) AddRecentFile(path string) {
	// Remove if already exists
	for i, file := range c.RecentFiles {
		if file == path {
			c.RecentFiles = append(c.RecentFiles[:i], c.RecentFiles[i+1:]...)
			break
		}
	}

	// Add to the front
	c.RecentFiles = append([]string{path}, c.RecentFiles...)

	// Limit to 10 recent files
	if len(c.RecentFiles) > 10 {
		c.RecentFiles = c.RecentFiles[:10]
	}

	// Save config
	c.Save()
}

// GetAutoSaveDelayDuration returns the autosave delay as a time.Duration
func (c *Config) GetAutoSaveDelayDuration() time.Duration {
	return time.Duration(c.AutoSaveDelay) * time.Second
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create config directory path
	configDir := filepath.Join(homeDir, ".markdown-editor")

	// Create config file path
	configPath := filepath.Join(configDir, "config.json")

	return configPath, nil
}
