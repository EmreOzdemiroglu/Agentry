package main

import (
	"fmt"

	"agentic-creator/internal/config"
	"agentic-creator/internal/models"
	"agentic-creator/internal/services"
	"agentic-creator/internal/ui"
)

// App holds all the application dependencies
type App struct {
	config           *config.Config
	inputReader      *ui.InputReader
	ollamaService    services.OllamaService
	chainService     services.ChainService
	executionService services.ExecutionService
	logService       services.LogService
}

// NewApp creates a new application instance
func NewApp() *App {
	cfg := config.Default()
	inputReader := ui.NewInputReader()
	ollamaService := services.NewOllamaService(cfg)
	chainService := services.NewChainService(cfg)
	executionService := services.NewExecutionService(ollamaService)
	logService := services.NewLogService(cfg)

	return &App{
		config:           cfg,
		inputReader:      inputReader,
		ollamaService:    ollamaService,
		chainService:     chainService,
		executionService: executionService,
		logService:       logService,
	}
}

func main() {
	_ = NewApp() // Initialize app (unused for now in this test version)
	fmt.Println("Welcome to the Agentic Structure Creator!")

	// Test creating basic structures
	agent := models.Agent{
		Name:         "test",
		Model:        "llama2",
		SystemPrompt: "You are helpful",
		UserPrompt:   "Help me",
	}
	
	chain := models.AgentChain{
		Name:   "test-chain",
		Agents: []models.Agent{agent},
		Loop:   false,
	}
	
	fmt.Printf("Created chain: %s with %d agents\n", chain.Name, len(chain.Agents))
	ui.DisplayChain(&chain)
}