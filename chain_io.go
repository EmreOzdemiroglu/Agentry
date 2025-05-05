package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	chainsDir = "Agentries"
	logsDir   = "chatruns"
)

// ensureDir checks if a directory exists, creates it if not.
func ensureDir(dirName string) error {
	// Check if the path exists and if it's a directory
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		fmt.Printf("Creating directory: %s\n", dirName)
		return os.MkdirAll(dirName, 0755) // 0755 standard permissions
	} else if err != nil {
		// Other error (e.g., permission denied)
		return fmt.Errorf("error checking directory %s: %w", dirName, err)
	} else if !info.IsDir() {
		// Path exists but is not a directory
		return fmt.Errorf("path %s exists but is not a directory", dirName)
	}
	// Directory exists
	return nil
}

// findSavedChains searches for .json files in the chains directory
func findSavedChains() ([]string, error) {
	pattern := filepath.Join(chainsDir, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err // Glob errors are usually about pattern syntax
	}

	// Check if the directory exists if no matches are found, to give a better message
	if len(matches) == 0 {
		_ = ensureDir(chainsDir) // Try to create it if it doesn't exist, ignore error here
	}

	return matches, nil
}

// loadChain reads a JSON file (expected full path including directory) and unmarshals it
func loadChain(fullPath string) (AgentChain, error) {
	var chain AgentChain
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return chain, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	err = json.Unmarshal(fileData, &chain)
	if err != nil {
		return chain, fmt.Errorf("failed to unmarshal JSON from %s: %w", fullPath, err)
	}

	// Optional: Validate the loaded chain structure if needed
	if chain.Name == "" || len(chain.Agents) == 0 {
		// fmt.Printf("Warning: Loaded chain '%s' might be incomplete.\n", fullPath)
	}

	return chain, nil
}

// saveChain marshals the AgentChain to JSON and saves it to the chains directory
func saveChain(chain AgentChain, baseFilename string) error {
	if err := ensureDir(chainsDir); err != nil {
		return fmt.Errorf("could not ensure directory %s exists: %w", chainsDir, err)
	}

	fullPath := filepath.Join(chainsDir, baseFilename)

	jsonData, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chain to JSON: %w", err)
	}

	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write chain to file '%s': %w", fullPath, err)
	}

	return nil // Success
}

// saveRunLog saves the captured log buffer to a timestamped markdown file in the logs directory
func saveRunLog(chain AgentChain, logBuffer *bytes.Buffer, startTime time.Time) error {
	if err := ensureDir(logsDir); err != nil {
		return fmt.Errorf("could not ensure directory %s exists: %w", logsDir, err)
	}

	// Sanitize chain name for filename
	sanitizedName := strings.ReplaceAll(chain.Name, " ", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "/", "-")
	// Add more sanitization if needed

	baseFilename := fmt.Sprintf("chatruns_%s_%s.md",
		startTime.Format("20060102_150405"),
		sanitizedName)

	fullPath := filepath.Join(logsDir, baseFilename)

	err := os.WriteFile(fullPath, logBuffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write log file '%s': %w", fullPath, err)
	}
	fmt.Printf("Run log saved to %s\n", fullPath) // Inform user of the full path
	return nil
} 