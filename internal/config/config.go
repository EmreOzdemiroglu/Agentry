package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	OllamaURL string
	ChainsDir string
	LogsDir   string
}

// Default returns a configuration with default values
func Default() *Config {
	return &Config{
		OllamaURL: "http://localhost:11434/api/generate",
		ChainsDir: "Agentries",
		LogsDir:   "chatruns",
	}
}

// EnsureDir checks if a directory exists, creates it if not
func EnsureDir(dirName string) error {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0755)
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return os.ErrExist
	}
	return nil
}

// GetChainPath returns the full path for a chain file
func (c *Config) GetChainPath(filename string) string {
	return filepath.Join(c.ChainsDir, filename)
}

// GetLogPath returns the full path for a log file
func (c *Config) GetLogPath(filename string) string {
	return filepath.Join(c.LogsDir, filename)
}