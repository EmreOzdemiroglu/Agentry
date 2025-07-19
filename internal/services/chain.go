package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"agentic-creator/internal/config"
	"agentic-creator/internal/models"
)

// ChainServiceImpl implements the ChainService interface
type ChainServiceImpl struct {
	config *config.Config
}

// NewChainService creates a new chain service
func NewChainService(cfg *config.Config) *ChainServiceImpl {
	return &ChainServiceImpl{
		config: cfg,
	}
}

// FindSavedChains searches for .json files in the chains directory
func (s *ChainServiceImpl) FindSavedChains() ([]string, error) {
	pattern := filepath.Join(s.config.ChainsDir, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find chain files: %w", err)
	}

	if len(matches) == 0 {
		_ = config.EnsureDir(s.config.ChainsDir)
	}

	return matches, nil
}

// LoadChain reads a JSON file and unmarshals it into an AgentChain
func (s *ChainServiceImpl) LoadChain(fullPath string) (*models.AgentChain, error) {
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	var chain models.AgentChain
	if err := json.Unmarshal(fileData, &chain); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", fullPath, err)
	}

	if err := s.ValidateChain(&chain); err != nil {
		return nil, fmt.Errorf("loaded chain is invalid: %w", err)
	}

	return &chain, nil
}

// SaveChain marshals the AgentChain to JSON and saves it to the chains directory
func (s *ChainServiceImpl) SaveChain(chain *models.AgentChain, filename string) error {
	if err := s.ValidateChain(chain); err != nil {
		return fmt.Errorf("cannot save invalid chain: %w", err)
	}

	if err := config.EnsureDir(s.config.ChainsDir); err != nil {
		return fmt.Errorf("could not ensure directory %s exists: %w", s.config.ChainsDir, err)
	}

	fullPath := s.config.GetChainPath(filename)

	jsonData, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chain to JSON: %w", err)
	}

	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write chain to file '%s': %w", fullPath, err)
	}

	return nil
}

// ValidateChain validates the structure and content of an AgentChain
func (s *ChainServiceImpl) ValidateChain(chain *models.AgentChain) error {
	if chain == nil {
		return ValidationError{Field: "chain", Message: "chain cannot be nil"}
	}

	if strings.TrimSpace(chain.Name) == "" {
		return ValidationError{Field: "name", Message: "chain name cannot be empty"}
	}

	if len(chain.Agents) == 0 {
		return ValidationError{Field: "agents", Message: "chain must have at least one agent"}
	}

	for i, agent := range chain.Agents {
		if err := s.validateAgent(&agent, i); err != nil {
			return err
		}
	}

	return nil
}

// validateAgent validates a single agent
func (s *ChainServiceImpl) validateAgent(agent *models.Agent, index int) error {
	prefix := fmt.Sprintf("agent[%d]", index)

	if strings.TrimSpace(agent.Name) == "" {
		return ValidationError{Field: prefix + ".name", Message: "agent name cannot be empty"}
	}

	if strings.TrimSpace(agent.Model) == "" {
		return ValidationError{Field: prefix + ".model", Message: "agent model cannot be empty"}
	}

	if strings.TrimSpace(agent.SystemPrompt) == "" {
		return ValidationError{Field: prefix + ".system_prompt", Message: "agent system prompt cannot be empty"}
	}

	if strings.TrimSpace(agent.UserPrompt) == "" {
		return ValidationError{Field: prefix + ".user_prompt", Message: "agent user prompt cannot be empty"}
	}

	return nil
}