package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func loadAndRunChain() {
	fmt.Println("\n--- Load Existing Chain ---")
	savedChains, err := findSavedChains()
	if err != nil {
		fmt.Printf("Error finding saved chains: %v\n", err)
		return
	}

	if len(savedChains) == 0 {
		fmt.Printf("No saved chains (.json files) found in the '%s' directory.\n", chainsDir)
		return
	}

	fmt.Println("Available Chains:")
	for i, fullPath := range savedChains {
		// Display without the directory and .json extension for cleaner look
		baseName := filepath.Base(fullPath)
		chainName := strings.TrimSuffix(baseName, ".json")
		fmt.Printf("  %d: %s (%s)\n", i+1, chainName, baseName)
	}

	// Select chain to load
	selectedIndex := -1
	for {
		fmt.Printf("Select a chain number to load (1-%d) or 0 to cancel: ", len(savedChains))
		input := readInput("")
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(savedChains) {
			fmt.Println("Invalid selection.")
			continue
		}
		if choice == 0 {
			fmt.Println("Load cancelled.")
			return // Cancelled
		}
		selectedIndex = choice - 1
		break
	}

	selectedFullPath := savedChains[selectedIndex]
	loadedChain, err := loadChain(selectedFullPath)
	if err != nil {
		fmt.Printf("Error loading chain from %s: %v\n", selectedFullPath, err)
		return
	}

	fmt.Printf("\nChain '%s' loaded successfully from %s\n", loadedChain.Name, selectedFullPath)
	fmt.Println("-------------------------------------")
	displayChain(loadedChain)

	// Ask what to do with the loaded chain
	for {
		fmt.Println("\nOptions for loaded chain:")
		fmt.Println("  [R]un this chain")
		fmt.Println("  [E]dit this chain")
		fmt.Println("  [B]ack to main menu")
		choice := strings.ToUpper(readInput("> "))

		switch choice {
		case "R":
			runChain(loadedChain)
			return // Return to main menu after running
		case "E":
			baseFilename := filepath.Base(selectedFullPath)
			edited := editChain(&loadedChain, baseFilename)
			if edited {
				fmt.Println("Chain updated. Returning to options for this chain.")
			} else {
				fmt.Println("Edits cancelled. Returning to options for this chain.")
			}
		case "B":
			return // Back to main menu
		default:
			fmt.Println("Invalid choice. Please enter R, E, or B.")
		}
	}
}

func createNewChain() {
	fmt.Println("\n--- Creating New Agent Chain ---")

	// Fetch models ONCE for the creation process
	models, err := listOllamaModels()
	if err != nil {
		fmt.Printf("Error listing Ollama models: %v\n", err)
		os.Exit(1) // Exit if models can't be listed initially
	}
	if len(models) == 0 {
		fmt.Println("No Ollama models found. Please ensure Ollama is running and models are installed.")
		os.Exit(1)
	}

	fmt.Println("\nAvailable Ollama Models:")
	for i, model := range models {
		fmt.Printf("%d: %s\n", i+1, model)
	}

	// --- Model selection logic ---
	selectedModelIndex := selectModel(len(models))
	selectedModel := models[selectedModelIndex]
	fmt.Printf("Selected model: %s\n", selectedModel)

	// --- Get system and user prompts ---
	fmt.Println("\nEnter the System Prompt for the first agent:")
	systemPrompt := readInput("> ")
	fmt.Println("\nEnter the User Prompt template for the first agent:")
	userPrompt := readInput("> ")

	// --- Get Agent Name ---
	fmt.Println("\nEnter a unique name for this first agent (e.g., 'Summarizer', 'Translator'):")
	firstAgentName := readInput("> ")
	for firstAgentName == "" {
		fmt.Println("Agent name cannot be empty. Please enter a name:")
		firstAgentName = readInput("> ")
	}

	// --- Create first agent ---
	firstAgent := Agent{
		Name:         firstAgentName,
		Model:        selectedModel,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}

	chain := AgentChain{
		Agents: []Agent{firstAgent},
	}

	fmt.Printf("\nAgent %d: '%s' ([%s]) added to the chain.\n", len(chain.Agents), firstAgent.Name, firstAgent.Model)
	fmt.Println("-------------------------------------")
	displayChain(chain)

	// --- Loop for adding more agents ---
	for {
		fmt.Println("\nOptions:")
		fmt.Println("  [A]dd another agent to the chain")
		fmt.Println("  [S]ave chain")
		choice := strings.ToUpper(readInput("> "))

		if choice == "A" {
			// Call the refactored function
			agentAdded := addAgentInteractive(&chain, models)
			if agentAdded {
				fmt.Println("-------------------------------------")
				displayChain(chain)
			}
		} else if choice == "S" {
			fmt.Println("\nEnter a name for this agent chain (e.g., 'my-translator'):")
			chainName := readInput("> ")
			if chainName == "" {
				fmt.Println("Save cancelled. Chain name cannot be empty.")
				continue // Go back to Add/Save prompt
			}
			chain.Name = chainName // Set the name in the struct
			baseFilename := chainName + ".json"

			// Ask if the chain should loop
			fmt.Println("\nShould this chain loop? (Output of last agent feeds back to first) [Y/N]")
			loopChoice := strings.ToUpper(readInput("> "))
			chain.Loop = (loopChoice == "Y")

			err := saveChain(chain, baseFilename)
			if err != nil {
				fmt.Printf("Error saving chain: %v\n", err)
				continue
			} else {
				fmt.Printf("Agent chain '%s' saved successfully to %s (Looping: %v)\n", chain.Name, baseFilename, chain.Loop)
				// Ask if user wants to run the newly saved chain
				fmt.Println("\nDo you want to run this chain now? [Y/N]")
				runChoice := strings.ToUpper(readInput("> "))
				if runChoice == "Y" {
					runChain(chain)
				}
				break // Exit the loop after successful save
			}
		} else {
			fmt.Println("Invalid choice. Please enter 'A' or 'S'.")
		}
	}

	fmt.Println("Exiting chain creation.")
}

// addAgentInteractive handles the interactive process of adding a new or reused agent
// It modifies the chain directly and returns true if an agent was added.
func addAgentInteractive(chain *AgentChain, models []string) bool {
	fmt.Println("\n--- Add Next Agent --- ")

	// Option to Create New or Reuse Existing
	fmt.Println("\nOptions:")
	fmt.Println("  [N]ew Agent")
	fmt.Println("  [R]euse Existing Agent from this chain")
	fmt.Println("  [B]ack")
	addChoice := strings.ToUpper(readInput("> "))

	var nextAgent Agent
	var agentAdded bool = false

	if addChoice == "N" {
		// --- Create NEW agent ---
		fmt.Println("\nAvailable Ollama Models:")
		for i, model := range models {
			fmt.Printf("%d: %s\n", i+1, model)
		}
		selectedModelIndex := selectModel(len(models))
		selectedModel := models[selectedModelIndex]
		fmt.Printf("Selected model: %s\n", selectedModel)

		fmt.Println("\nEnter the System Prompt for this new agent:")
		systemPrompt := readInput("> ")
		fmt.Println("\nEnter the User Prompt template for this new agent:")
		userPrompt := readInput("> ")

		fmt.Println("\nEnter a unique name for this new agent:")
		newAgentName := readInput("> ")
		// TODO: Add validation to ensure name is unique within the chain
		for newAgentName == "" || agentNameExists(chain, newAgentName) { // Basic & uniqueness validation
			if newAgentName == "" {
				fmt.Println("Agent name cannot be empty.")
			} else {
				fmt.Printf("Agent name '%s' already exists. Please enter a unique name:\n", newAgentName)
			}
			newAgentName = readInput("> ")
		}

		nextAgent = Agent{
			Name:         newAgentName,
			Model:        selectedModel,
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt,
		}
		agentAdded = true

	} else if addChoice == "R" {
		// --- Reuse EXISTING agent ---
		if len(chain.Agents) == 0 {
			fmt.Println("No existing agents in the chain to reuse.")
			return false // Go back
		}
		fmt.Println("\nSelect an agent definition to reuse:")
		for i, agent := range chain.Agents {
			fmt.Printf("  %d: %s (%s)\n", i+1, agent.Name, agent.Model)
		}

		selectedIndex := -1
		for {
			fmt.Printf("Enter the number of the agent to reuse (1-%d): ", len(chain.Agents))
			input := readInput("")
			choiceNum, err := strconv.Atoi(input)
			if err != nil || choiceNum < 1 || choiceNum > len(chain.Agents) {
				fmt.Println("Invalid selection.")
				continue
			}
			selectedIndex = choiceNum - 1
			break
		}
		baseAgent := chain.Agents[selectedIndex]
		fmt.Printf("Reusing agent definition '%s' (%s)\n", baseAgent.Name, baseAgent.Model)

		// Ask about User Prompt
		fmt.Println("\nDo you want to: [K]eep the original user prompt OR [M]odify it for this step?")
		promptChoice := strings.ToUpper(readInput("> "))
		newUserPrompt := baseAgent.UserPrompt // Default to original
		if promptChoice == "M" {
			fmt.Println("\nEnter the new User Prompt template for this reused agent instance:")
			newUserPrompt = readInput("> ")
		}

		// Get a unique name for this *instance*
		fmt.Printf("\nEnter a unique name for this instance of the reused agent (e.g., '%s_step%d'):\n", baseAgent.Name, len(chain.Agents)+1)
		reusedAgentName := readInput("> ")
		// TODO: Add better unique name suggestion/validation
		for reusedAgentName == "" || agentNameExists(chain, reusedAgentName) { // Basic & uniqueness validation
			if reusedAgentName == "" {
				fmt.Println("Agent name cannot be empty.")
			} else {
				fmt.Printf("Agent name '%s' already exists. Please enter a unique name:\n", reusedAgentName)
			}
			reusedAgentName = readInput("> ")
		}

		nextAgent = Agent{
			Name:         reusedAgentName,
			Model:        baseAgent.Model,
			SystemPrompt: baseAgent.SystemPrompt,
			UserPrompt:   newUserPrompt,
		}
		agentAdded = true

	} else if addChoice == "B" {
		return false // Go back without adding
	} else {
		fmt.Println("Invalid choice. Please enter 'N', 'R', or 'B'.")
		return false // Go back without adding
	}

	// Append the newly created or reused agent if one was successfully defined
	if agentAdded {
		chain.Agents = append(chain.Agents, nextAgent)
		fmt.Printf("\nAgent %d: '%s' ([%s]) added to the chain.\n", len(chain.Agents), nextAgent.Name, nextAgent.Model)
		return true
	}
	return false // Should not happen if logic is correct, but safety return
}

// agentNameExists checks if an agent name already exists in the chain
func agentNameExists(chain *AgentChain, name string) bool {
	for _, agent := range chain.Agents {
		if agent.Name == name {
			return true
		}
	}
	return false
}

// runChain executes the agent chain with initial input
func runChain(chain AgentChain) {
	var logBuffer bytes.Buffer // Buffer to store the log for saving
	startTime := time.Now()

	// --- Log Header ---
	logHeader := fmt.Sprintf("# Agent Chain Run Log\n\n**Chain Name:** %s\n**Looping:** %v\n**Started:** %s\n\n",
		chain.Name, chain.Loop, startTime.Format(time.RFC1123))
	logBuffer.WriteString(logHeader)
	fmt.Print(logHeader) // Also print header to console

	fmt.Println("Enter the initial input for the chain:")
	initialInput := readInput("> ")
	logBuffer.WriteString(fmt.Sprintf("## Initial Input\n\n```\n%s\n```\n\n", initialInput))
	fmt.Printf("Initial Input recorded.\n")

	// Ollama API endpoint
	ollamaURL := "http://localhost:11434/api/generate"

	loopInput := initialInput
	loopCounter := 0

	for {
		loopCounter++
		iterationStartTime := time.Now()
		currentInput := loopInput
		var currentContext []int

		iterHeader := fmt.Sprintf("## Loop Iteration %d\n\n**Start Time:** %s\n**Iteration Input:**\n\n```\n%s\n```\n\n",
			loopCounter, iterationStartTime.Format(time.RFC1123), currentInput)
		logBuffer.WriteString(iterHeader)
		fmt.Printf("\n--- Starting Loop Iteration %d ---\n", loopCounter)
		fmt.Printf("Iteration Input:\n%s\n", currentInput)
		fmt.Println("-------------------------------------")

		for _, agent := range chain.Agents {
			agentStartTime := time.Now()
			agentLog := fmt.Sprintf("### Agent: %s (%s)\n\n", agent.Name, agent.Model)
			logBuffer.WriteString(agentLog)
			fmt.Printf("\n[%d: %s (%s)] Processing...\n", loopCounter, agent.Name, agent.Model)

			// Combine prompt
			combinedPrompt := currentInput + "\n\n" + agent.UserPrompt
			logBuffer.WriteString(fmt.Sprintf("**System Prompt:**\n```\n%s\n```\n", agent.SystemPrompt))
			logBuffer.WriteString(fmt.Sprintf("**Combined User Prompt (Input + Template):**\n```\n%s\n```\n", combinedPrompt))

			// Prepare request
			requestPayload := OllamaRequest{
				Model:   agent.Model,
				System:  agent.SystemPrompt,
				Prompt:  combinedPrompt,
				Stream:  false,
				Context: currentContext,
			}

			jsonData, err := json.Marshal(requestPayload)
			if err != nil {
				errMsg := fmt.Sprintf("**ERROR Marshaling Request:** %v\n", err)
				logBuffer.WriteString(errMsg)
				fmt.Printf(errMsg)
				return
			}

			// Make request and handle response
			req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(jsonData))
			if err != nil {
				errMsg := fmt.Sprintf("**ERROR Creating Request:** %v\n", err)
				logBuffer.WriteString(errMsg)
				fmt.Printf(errMsg)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errMsg := fmt.Sprintf("**ERROR Sending Request:** %v\nIs Ollama running at %s?\n", err, ollamaURL)
				logBuffer.WriteString(errMsg)
				fmt.Printf(errMsg)
				return
			}

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				errMsg := fmt.Sprintf("**ERROR Response from Ollama (Status %d):**\n```\n%s\n```\n", resp.StatusCode, string(bodyBytes))
				logBuffer.WriteString(errMsg)
				fmt.Printf(errMsg)
				return
			}

			// Decode response
			var ollamaResp OllamaResponse
			if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
				resp.Body.Close()
				errMsg := fmt.Sprintf("**ERROR Decoding Ollama Response:** %v\n", err)
				logBuffer.WriteString(errMsg)
				fmt.Printf(errMsg)
				return
			}
			resp.Body.Close()

			agentEndTime := time.Now()
			agentDuration := agentEndTime.Sub(agentStartTime)

			// Log response and update context/input
			currentInput = ollamaResp.Response
			currentContext = ollamaResp.Context

			agentResultLog := fmt.Sprintf("**Response (Duration: %s):**\n\n```\n%s\n```\n\n",
				agentDuration.Round(time.Millisecond), ollamaResp.Response)
			logBuffer.WriteString(agentResultLog)
			fmt.Printf("[%d: %s (%s)] Response (Duration: %s):\n%s\n",
				loopCounter, agent.Name, agent.Model, agentDuration.Round(time.Millisecond), ollamaResp.Response)
			fmt.Println("-------------------------------------")
		} // End of inner agent loop

		// --- Iteration finished ---
		iterationEndTime := time.Now()
		iterationDuration := iterationEndTime.Sub(iterationStartTime)

		iterSummary := fmt.Sprintf("**Finished Loop Iteration %d (Duration: %s)**\n\n**Iteration Output:**\n\n```\n%s\n```\n\n---\n\n",
			loopCounter, iterationDuration.Round(time.Millisecond), currentInput)
		logBuffer.WriteString(iterSummary)
		fmt.Printf("\n--- Finished Loop Iteration %d (Duration: %s) ---\n", loopCounter, iterationDuration.Round(time.Millisecond))
		fmt.Printf("Iteration Output:\n%s\n", currentInput)
		fmt.Println("=====================================")

		// Update loopInput for the NEXT iteration
		loopInput = currentInput

		// Check if we should continue looping
		if !chain.Loop {
			break // Not a looping chain, exit outer loop after one iteration
		}

		// Ask user if they want to continue the loop
		fmt.Println("\nRun another loop iteration? [Y/N]")
		continueChoice := strings.ToUpper(readInput("> "))
		if continueChoice != "Y" {
			break // Exit the outer loop
		}
	} // End of outer loop

	// --- Chain finished ---
	finalMsg := fmt.Sprintf("## Chain Execution Finished\n\n**Final Output:**\n\n```\n%s\n```\n", loopInput)
	logBuffer.WriteString(finalMsg)
	fmt.Println("\n--- Chain execution finished --- Final Output ---")
	fmt.Printf("%s\n", loopInput)
	fmt.Println("===============================================")

	// Ask to save log
	fmt.Println("\nDo you want to save the full run log to a Markdown file? [Y/N]")
	saveChoice := strings.ToUpper(readInput("> "))
	if saveChoice == "Y" {
		err := saveRunLog(chain, &logBuffer, startTime)
		if err != nil {
			fmt.Printf("Error saving run log: %v\n", err)
		} else {
			fmt.Println("Run log saved successfully.")
		}
	}
}

// editChain provides an interface to modify a loaded chain
// Returns true if changes were saved, false otherwise (cancelled).
// originalBaseFilename is just the file name part (e.g., my-chain.json)
func editChain(chain *AgentChain, originalBaseFilename string) bool {
	fmt.Printf("\n--- Editing Chain: %s ---\n", chain.Name)

	// Fetch models needed for adding agents
	models, err := listOllamaModels()
	if err != nil {
		fmt.Printf("Error listing Ollama models (needed for editing): %v\n", err)
		fmt.Println("Cannot add new agents without model list. Returning to previous menu.")
		return false // Cannot proceed with editing involving adding
	}
	if len(models) == 0 {
		fmt.Println("No Ollama models found. Cannot add new agents.")
		// Allow deleting existing agents, but not adding
	}

	for {
		fmt.Println("\nCurrent structure:")
		displayChain(*chain)
		fmt.Println("Edit Options:")
		fmt.Println("  [A]dd Agent")
		fmt.Println("  [D]elete Agent")
		fmt.Println("  [S]ave Changes")
		fmt.Println("  [C]ancel Edits")
		choice := strings.ToUpper(readInput("> "))

		switch choice {
		case "A":
			if len(models) == 0 {
				fmt.Println("Cannot add agent: No Ollama models available.")
				continue
			}
			addAgentInteractive(chain, models) // Modifies chain directly
			// Display happens at the start of the next loop iteration
		case "D":
			if len(chain.Agents) == 0 {
				fmt.Println("No agents to delete.")
				continue
			}
			fmt.Println("Select agent number to delete:")
			delIndex := -1
			for {
				fmt.Printf("Enter number (1-%d): ", len(chain.Agents))
				input := readInput("")
				choiceNum, err := strconv.Atoi(input)
				if err != nil || choiceNum < 1 || choiceNum > len(chain.Agents) {
					fmt.Println("Invalid selection.")
					continue
				}
				delIndex = choiceNum - 1
				break
			}
			deletedAgentName := chain.Agents[delIndex].Name
			// Slice trick to delete element at delIndex
			chain.Agents = append(chain.Agents[:delIndex], chain.Agents[delIndex+1:]...)
			fmt.Printf("Agent '%s' deleted.\n", deletedAgentName)
		case "S":
			// Ask about loop flag again when saving edits?
			// For simplicity, let's keep the original loop flag unless explicitly changed.
			// Or we can ask:
			fmt.Println("\nShould this chain loop? (Current: %v) [Y/N/K(eep)]", chain.Loop)
			loopChoice := strings.ToUpper(readInput("> "))
			if loopChoice == "Y" {
				chain.Loop = true
			} else if loopChoice == "N" {
				chain.Loop = false
			} // Else K or invalid, keep original

			err := saveChain(*chain, originalBaseFilename)
			if err != nil {
				fmt.Printf("Error saving changes: %v\n", err)
				// Stay in edit mode
			} else {
				savedPath := filepath.Join(chainsDir, originalBaseFilename)
				fmt.Printf("Changes saved successfully to %s\n", savedPath)
				return true // Indicate changes were saved
			}
		case "C":
			fmt.Println("Cancelling edits. No changes saved.")
			return false // Indicate changes were cancelled
		default:
			fmt.Println("Invalid choice. Please enter A, D, S, or C.")
		}
	}
} 