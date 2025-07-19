package services

import (
	"bytes"
	"fmt"
	"time"

	"agentic-creator/internal/models"
)

// ExecutionServiceImpl implements the ExecutionService interface
type ExecutionServiceImpl struct {
	ollamaService OllamaService
}

// NewExecutionService creates a new execution service
func NewExecutionService(ollamaService OllamaService) *ExecutionServiceImpl {
	return &ExecutionServiceImpl{
		ollamaService: ollamaService,
	}
}

// RunChain executes the agent chain with initial input
func (s *ExecutionServiceImpl) RunChain(chain *models.AgentChain, initialInput string) (*ExecutionResult, error) {
	var logBuffer bytes.Buffer
	startTime := time.Now()

	s.writeLogHeader(&logBuffer, chain, startTime)
	s.writeInitialInput(&logBuffer, initialInput)

	loopInput := initialInput
	loopCounter := 0

	for {
		loopCounter++
		iterationStartTime := time.Now()

		result, err := s.runIteration(chain, loopInput, loopCounter, &logBuffer, iterationStartTime)
		if err != nil {
			return nil, fmt.Errorf("iteration %d failed: %w", loopCounter, err)
		}

		loopInput = result

		if !chain.Loop {
			break
		}

		// In a real CLI, you'd ask the user here if they want to continue
		// For now, we'll just run once for non-looping chains
		break
	}

	duration := time.Since(startTime)
	s.writeFinalOutput(&logBuffer, loopInput)

	return &ExecutionResult{
		FinalOutput: loopInput,
		LogBuffer:   &logBuffer,
		StartTime:   startTime,
		Duration:    duration,
	}, nil
}

// runIteration runs a single iteration of the chain
func (s *ExecutionServiceImpl) runIteration(chain *models.AgentChain, input string, iteration int, logBuffer *bytes.Buffer, startTime time.Time) (string, error) {
	s.writeIterationHeader(logBuffer, iteration, startTime, input)

	currentInput := input
	var currentContext []int

	for _, agent := range chain.Agents {
		agentStartTime := time.Now()
		s.writeAgentHeader(logBuffer, iteration, &agent)

		result, context, err := s.runAgent(&agent, currentInput, currentContext, logBuffer)
		if err != nil {
			return "", fmt.Errorf("agent %s failed: %w", agent.Name, err)
		}

		agentDuration := time.Since(agentStartTime)
		s.writeAgentResult(logBuffer, iteration, &agent, result, agentDuration)

		currentInput = result
		currentContext = context
	}

	iterationDuration := time.Since(startTime)
	s.writeIterationResult(logBuffer, iteration, currentInput, iterationDuration)

	return currentInput, nil
}

// runAgent executes a single agent
func (s *ExecutionServiceImpl) runAgent(agent *models.Agent, input string, context []int, logBuffer *bytes.Buffer) (string, []int, error) {
	combinedPrompt := input + "\n\n" + agent.UserPrompt
	s.writeAgentPrompts(logBuffer, agent, combinedPrompt)

	request := models.OllamaRequest{
		Model:   agent.Model,
		System:  agent.SystemPrompt,
		Prompt:  combinedPrompt,
		Stream:  false,
		Context: context,
	}

	response, err := s.ollamaService.GenerateResponse(request)
	if err != nil {
		return "", nil, fmt.Errorf("ollama request failed: %w", err)
	}

	return response.Response, response.Context, nil
}

// Logging helper methods
func (s *ExecutionServiceImpl) writeLogHeader(logBuffer *bytes.Buffer, chain *models.AgentChain, startTime time.Time) {
	header := fmt.Sprintf("# Agent Chain Run Log\n\n**Chain Name:** %s\n**Looping:** %v\n**Started:** %s\n\n",
		chain.Name, chain.Loop, startTime.Format(time.RFC1123))
	logBuffer.WriteString(header)
}

func (s *ExecutionServiceImpl) writeInitialInput(logBuffer *bytes.Buffer, input string) {
	logBuffer.WriteString(fmt.Sprintf("## Initial Input\n\n```\n%s\n```\n\n", input))
}

func (s *ExecutionServiceImpl) writeIterationHeader(logBuffer *bytes.Buffer, iteration int, startTime time.Time, input string) {
	header := fmt.Sprintf("## Loop Iteration %d\n\n**Start Time:** %s\n**Iteration Input:**\n\n```\n%s\n```\n\n",
		iteration, startTime.Format(time.RFC1123), input)
	logBuffer.WriteString(header)
}

func (s *ExecutionServiceImpl) writeAgentHeader(logBuffer *bytes.Buffer, iteration int, agent *models.Agent) {
	agentLog := fmt.Sprintf("### Agent: %s (%s)\n\n", agent.Name, agent.Model)
	logBuffer.WriteString(agentLog)
}

func (s *ExecutionServiceImpl) writeAgentPrompts(logBuffer *bytes.Buffer, agent *models.Agent, combinedPrompt string) {
	logBuffer.WriteString(fmt.Sprintf("**System Prompt:**\n```\n%s\n```\n", agent.SystemPrompt))
	logBuffer.WriteString(fmt.Sprintf("**Combined User Prompt (Input + Template):**\n```\n%s\n```\n", combinedPrompt))
}

func (s *ExecutionServiceImpl) writeAgentResult(logBuffer *bytes.Buffer, iteration int, agent *models.Agent, result string, duration time.Duration) {
	agentResultLog := fmt.Sprintf("**Response (Duration: %s):**\n\n```\n%s\n```\n\n",
		duration.Round(time.Millisecond), result)
	logBuffer.WriteString(agentResultLog)
}

func (s *ExecutionServiceImpl) writeIterationResult(logBuffer *bytes.Buffer, iteration int, output string, duration time.Duration) {
	iterSummary := fmt.Sprintf("**Finished Loop Iteration %d (Duration: %s)**\n\n**Iteration Output:**\n\n```\n%s\n```\n\n---\n\n",
		iteration, duration.Round(time.Millisecond), output)
	logBuffer.WriteString(iterSummary)
}

func (s *ExecutionServiceImpl) writeFinalOutput(logBuffer *bytes.Buffer, output string) {
	finalMsg := fmt.Sprintf("## Chain Execution Finished\n\n**Final Output:**\n\n```\n%s\n```\n", output)
	logBuffer.WriteString(finalMsg)
}