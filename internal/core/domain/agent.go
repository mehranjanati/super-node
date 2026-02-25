package domain

import "time"

// AgentType defines the type of agent
type AgentType string

const (
	AgentTypeTrading   AgentType = "trading"
	AgentTypeAnalytics AgentType = "analytics"
	AgentTypeSocial    AgentType = "social"
	AgentTypeContent   AgentType = "content"
	AgentTypeCustom    AgentType = "custom"
)

// AgentStatus defines the current status of the agent
type AgentStatus string

const (
	AgentStatusActive    AgentStatus = "active"
	AgentStatusPaused    AgentStatus = "paused"
	AgentStatusError     AgentStatus = "error"
	AgentStatusDeploying AgentStatus = "deploying"
)

// AgentPerformance metrics for the agent
type AgentPerformance struct {
	ROI         float64   `json:"roi"`          // Return on Investment percentage
	Trades      int       `json:"trades"`       // Total number of trades/actions
	Uptime      float64   `json:"uptime"`       // Uptime percentage
	SuccessRate float64   `json:"successRate"`  // Success rate percentage
	LastActive  time.Time `json:"lastActive"`
}

// Agent represents an AI agent in the system
type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        AgentType              `json:"type"`
	Status      AgentStatus            `json:"status"`
	Performance AgentPerformance       `json:"performance"`
	OwnerID     string                 `json:"owner"`
	Avatar      string                 `json:"avatar,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// AgentConfig defines the dynamic personality and capabilities of an agent
// This is now a subset or part of the Agent's Config map
type AgentConfig struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	SystemPrompt   string   `json:"system_prompt"`   // The core personality/instruction
	ModelID        string   `json:"model_id"`        // e.g., "deepseek-r1-distill-llama-8b"
	LoRAAdapterID  string   `json:"lora_adapter_id"` // Path to specific fine-tuned weights
	Temperature    float64  `json:"temperature"`
	Tools          []string `json:"tools"`           // Enabled MCP tools
	RewardCriteria string   `json:"reward_criteria"` // Description of what constitutes "success"
}

// TrainingSample represents a single interaction for GRPO/RLHF training
type TrainingSample struct {
	Timestamp      int64       `json:"timestamp"`
	AgentID        string      `json:"agent_id"`
	InputContext   interface{} `json:"input_context"`   // Market data, news
	ReasoningTrace string      `json:"reasoning_trace"` // The "Thought Stream"
	Action         string      `json:"action"`          // The decision (BUY/SELL)
	Outcome        string      `json:"outcome"`         // User Approved? Profit?
	RewardScore    float64     `json:"reward_score"`    // +1.0 for approval, -1.0 for rejection
}

// MarketAnalysisResult represents the output from the Wasm/Rust agent
type MarketAnalysisResult struct {
	Strategy  string `json:"strategy"`
	TopPick   string `json:"top_pick"`
	RiskLevel string `json:"risk_level"`
	Reasoning string `json:"reasoning"`
}
