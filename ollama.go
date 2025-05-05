package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// listOllamaModels executes `ollama list` and returns a slice of model names
func listOllamaModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		// Check if the error is because ollama command is not found
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an error code
			return nil, fmt.Errorf("'ollama list' failed. Is Ollama running and installed? %v, stderr: %s", err, string(exitErr.Stderr))
		}
		// Command not found or other errors
		return nil, fmt.Errorf("failed to execute 'ollama list': %w. Is Ollama installed and in PATH?", err)

	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var models []string
	headerSkipped := false

	for scanner.Scan() {
		line := scanner.Text()
		if !headerSkipped {
			// Assuming the first line is the header NAME ID SIZE MODIFIED
			if strings.HasPrefix(strings.ToUpper(line), "NAME") {
				headerSkipped = true
				continue
			}
			// Handle case where there might not be a header or format is unexpected
			// For now, we'll try processing anyway, might need better parsing
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0]) // Assuming the model name is the first column
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ollama list output: %w", err)
	}

	return models, nil
} 