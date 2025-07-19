package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"agentic-creator/internal/config"
	"agentic-creator/internal/models"
)

// OllamaServiceImpl implements the OllamaService interface
type OllamaServiceImpl struct {
	config *config.Config
	client *http.Client
}

// NewOllamaService creates a new Ollama service
func NewOllamaService(cfg *config.Config) *OllamaServiceImpl {
	return &OllamaServiceImpl{
		config: cfg,
		client: &http.Client{},
	}
}

// ListModels executes `ollama list` and returns a slice of model names
func (s *OllamaServiceImpl) ListModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("'ollama list' failed. Is Ollama running and installed? %v, stderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute 'ollama list': %w. Is Ollama installed and in PATH?", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var models []string
	headerSkipped := false

	for scanner.Scan() {
		line := scanner.Text()
		if !headerSkipped {
			if strings.HasPrefix(strings.ToUpper(line), "NAME") {
				headerSkipped = true
				continue
			}
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ollama list output: %w", err)
	}

	return models, nil
}

// GenerateResponse sends a request to Ollama and returns the response
func (s *OllamaServiceImpl) GenerateResponse(req models.OllamaRequest) (*models.OllamaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.config.OllamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama at %s: %w", s.config.OllamaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	var ollamaResp models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return &ollamaResp, nil
}