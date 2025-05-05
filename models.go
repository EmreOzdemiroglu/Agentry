package main

// Agent defines the structure for a single agent in the chain
type Agent struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
	UserPrompt   string `json:"user_prompt"` // This is the template
}

// AgentChain represents the sequence of agents
type AgentChain struct {
	Name   string  `json:"name"`
	Agents []Agent `json:"agents"`
	Loop   bool    `json:"loop,omitempty"` // Added loop flag
}

// Ollama API Request Payload
type OllamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	System  string `json:"system"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"` // To hold context between calls if needed (optional for now)
}

// Ollama API Response Payload (when stream: false)
type OllamaResponse struct {
	Model             string    `json:"model"`
	CreatedAt         string `json:"created_at"`
	Response          string    `json:"response"`
	Done              bool      `json:"done"`
	Context           []int     `json:"context"` // Context to pass back for subsequent requests
	TotalDuration     int64     `json:"total_duration"`
	LoadDuration      int64     `json:"load_duration"`
	PromptEvalCount   int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount         int       `json:"eval_count"`
	EvalDuration      int64     `json:"eval_duration"`
} 