package services

import (
	"bytes"
	"time"

	"agentic-creator/internal/models"
)

// OllamaService defines the interface for Ollama operations
type OllamaService interface {
	ListModels() ([]string, error)
	GenerateResponse(req models.OllamaRequest) (*models.OllamaResponse, error)
}

// ChainService defines the interface for chain operations
type ChainService interface {
	FindSavedChains() ([]string, error)
	LoadChain(fullPath string) (*models.AgentChain, error)
	SaveChain(chain *models.AgentChain, filename string) error
	ValidateChain(chain *models.AgentChain) error
}

// ExecutionService defines the interface for chain execution
type ExecutionService interface {
	RunChain(chain *models.AgentChain, initialInput string) (*ExecutionResult, error)
}

// LogService defines the interface for logging operations
type LogService interface {
	SaveRunLog(chain *models.AgentChain, logBuffer *bytes.Buffer, startTime time.Time) error
}

// ExecutionResult holds the result of chain execution
type ExecutionResult struct {
	FinalOutput string
	LogBuffer   *bytes.Buffer
	StartTime   time.Time
	Duration    time.Duration
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}