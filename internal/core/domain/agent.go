package domain

// AgentConfig defines the dynamic personality and capabilities of an agent
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
