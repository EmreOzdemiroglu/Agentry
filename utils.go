package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// selectModel prompts the user to select a model by number
func selectModel(numModels int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Select a model number (1-%d): ", numModels)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > numModels {
			fmt.Println("Invalid selection. Please enter a number from the list.")
			continue
		}
		return choice - 1 // Return 0-based index
	}
}

// readInput reads a line of text from stdin after printing a prompt
func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// displayChain shows the current agent chain structure
func displayChain(chain AgentChain) {
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