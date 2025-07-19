package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"agentic-creator/internal/models"
)

// DisplayChain shows the current agent chain structure
func DisplayChain(chain *models.AgentChain) {
	fmt.Println("Current Chain Structure:")
	if len(chain.Agents) == 0 {
		fmt.Println("(empty)")
		return
	}
	for i, agent := range chain.Agents {
		if i > 0 {
			fmt.Println("   |")
			fmt.Println("   V")
		}
		fmt.Printf("[%d: %s (%s)]\n", i+1, agent.Name, agent.Model)
	}
	fmt.Println("-------------------------------------")
}

// DisplayModels shows a numbered list of available models
func DisplayModels(models []string) {
	fmt.Println("\nAvailable Ollama Models:")
	for i, model := range models {
		fmt.Printf("%d: %s\n", i+1, model)
	}
}

// DisplayChainList shows a numbered list of saved chains
func DisplayChainList(savedChains []string) {
	fmt.Println("Available Chains:")
	for i, fullPath := range savedChains {
		baseName := filepath.Base(fullPath)
		chainName := strings.TrimSuffix(baseName, ".json")
		fmt.Printf("  %d: %s (%s)\n", i+1, chainName, baseName)
	}
}

// DisplayAgentList shows a numbered list of agents for selection
func DisplayAgentList(agents []models.Agent) {
	fmt.Println("\nSelect an agent definition to reuse:")
	for i, agent := range agents {
		fmt.Printf("  %d: %s (%s)\n", i+1, agent.Name, agent.Model)
	}
}

// DisplayMainMenu shows the main application menu
func DisplayMainMenu() {
	fmt.Println("\n--- Main Menu ---")
	fmt.Println("  [C]reate New Agent Chain")
	fmt.Println("  [L]oad Existing Agent Chain")
	fmt.Println("  [Q]uit")
}

// DisplayChainOptions shows options for a loaded chain
func DisplayChainOptions() {
	fmt.Println("\nOptions for loaded chain:")
	fmt.Println("  [R]un this chain")
	fmt.Println("  [E]dit this chain")
	fmt.Println("  [B]ack to main menu")
}

// DisplayEditOptions shows options for editing a chain
func DisplayEditOptions() {
	fmt.Println("Edit Options:")
	fmt.Println("  [A]dd Agent")
	fmt.Println("  [D]elete Agent")
	fmt.Println("  [S]ave Changes")
	fmt.Println("  [C]ancel Edits")
}

// DisplayAddAgentOptions shows options for adding an agent
func DisplayAddAgentOptions() {
	fmt.Println("\nOptions:")
	fmt.Println("  [N]ew Agent")
	fmt.Println("  [R]euse Existing Agent from this chain")
	fmt.Println("  [B]ack")
}

// DisplayCreateChainOptions shows options during chain creation
func DisplayCreateChainOptions() {
	fmt.Println("\nOptions:")
	fmt.Println("  [A]dd another agent to the chain")
	fmt.Println("  [S]ave chain")
}

// PrintSeparator prints a separator line
func PrintSeparator() {
	fmt.Println("-------------------------------------")
}

// PrintDoubleSeparator prints a double separator line
func PrintDoubleSeparator() {
	fmt.Println("=====================================")
}