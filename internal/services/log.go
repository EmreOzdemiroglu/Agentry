package services

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"agentic-creator/internal/config"
	"agentic-creator/internal/models"
)

// LogServiceImpl implements the LogService interface
type LogServiceImpl struct {
	config *config.Config
}

// NewLogService creates a new log service
func NewLogService(cfg *config.Config) *LogServiceImpl {
	return &LogServiceImpl{
		config: cfg,
	}
}

// SaveRunLog saves the captured log buffer to a timestamped markdown file in the logs directory
func (s *LogServiceImpl) SaveRunLog(chain *models.AgentChain, logBuffer *bytes.Buffer, startTime time.Time) error {
	if err := config.EnsureDir(s.config.LogsDir); err != nil {
		return fmt.Errorf("could not ensure directory %s exists: %w", s.config.LogsDir, err)
	}

	sanitizedName := s.sanitizeFilename(chain.Name)
	baseFilename := fmt.Sprintf("chatruns_%s_%s.md",
		startTime.Format("20060102_150405"),
		sanitizedName)

	fullPath := s.config.GetLogPath(baseFilename)

	if err := os.WriteFile(fullPath, logBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write log file '%s': %w", fullPath, err)
	}

	return nil
}

// sanitizeFilename removes or replaces characters that are not suitable for filenames
func (s *LogServiceImpl) sanitizeFilename(name string) string {
	sanitized := strings.ReplaceAll(name, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, "\\", "-")
	sanitized = strings.ReplaceAll(sanitized, ":", "-")
	sanitized = strings.ReplaceAll(sanitized, "*", "-")
	sanitized = strings.ReplaceAll(sanitized, "?", "-")
	sanitized = strings.ReplaceAll(sanitized, "\"", "-")
	sanitized = strings.ReplaceAll(sanitized, "<", "-")
	sanitized = strings.ReplaceAll(sanitized, ">", "-")
	sanitized = strings.ReplaceAll(sanitized, "|", "-")
	return sanitized
}