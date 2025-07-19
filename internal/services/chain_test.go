package services

import (
	"testing"

	"agentic-creator/internal/config"
	"agentic-creator/internal/models"
)

func TestChainValidation(t *testing.T) {
	cfg := config.Default()
	service := NewChainService(cfg)

	tests := []struct {
		name    string
		chain   *models.AgentChain
		wantErr bool
	}{
		{
			name:    "nil chain",
			chain:   nil,
			wantErr: true,
		},
		{
			name: "empty name",
			chain: &models.AgentChain{
				Name:   "",
				Agents: []models.Agent{{Name: "test", Model: "test", SystemPrompt: "test", UserPrompt: "test"}},
			},
			wantErr: true,
		},
		{
			name: "no agents",
			chain: &models.AgentChain{
				Name:   "test",
				Agents: []models.Agent{},
			},
			wantErr: true,
		},
		{
			name: "valid chain",
			chain: &models.AgentChain{
				Name: "test-chain",
				Agents: []models.Agent{
					{Name: "agent1", Model: "llama2", SystemPrompt: "You are helpful", UserPrompt: "Help me"},
				},
			},
			wantErr: false,
		},
		{
			name: "agent with empty name",
			chain: &models.AgentChain{
				Name: "test-chain",
				Agents: []models.Agent{
					{Name: "", Model: "llama2", SystemPrompt: "You are helpful", UserPrompt: "Help me"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateChain(tt.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}