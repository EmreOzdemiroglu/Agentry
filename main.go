package main

import (
	"fmt"
	// "os" // Removed as it's no longer used directly in main.go
	"strings"
)

// --- Structs moved to models.go ---

func main() {
	fmt.Println("Welcome to the Agentic Structure Creator!")

	for {
		fmt.Println("\n--- Main Menu ---")
		fmt.Println("  [C]reate New Agent Chain")
		fmt.Println("  [L]oad Existing Agent Chain")
		fmt.Println("  [Q]uit")
		choice := strings.ToUpper(readInput("> "))

		switch choice {
		case "C":
			createNewChain() // This function now implicitly asks to run after saving
		case "L":
			loadAndRunChain()
		case "Q":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice. Please enter C, L, or Q.")
		}
	}
}

// --- Functions moved to chain_ops.go ---
// loadAndRunChain
// createNewChain
// addAgentInteractive
// agentNameExists
// runChain
// editChain

// --- Functions moved to chain_io.go ---
// findSavedChains
// loadChain
// saveChain
// saveRunLog

// --- Functions moved to ollama.go ---
// listOllamaModels

// --- Functions moved to utils.go ---
// selectModel
// readInput
// displayChain

// --- Potentially add other helper functions later --- 