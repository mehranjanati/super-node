package mlops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"nexus-super-node-v3/internal/core/domain"
)

// Collector handles the accumulation of training data for Unsloth/GRPO
type Collector struct {
	DataDir string
	mu      sync.Mutex
}

// NewCollector creates a new MLOps data collector
func NewCollector(dataDir string) *Collector {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create mlops data dir: %v\n", err)
	}
	return &Collector{
		DataDir: dataDir,
	}
}

// LogInteraction saves a training sample to a JSONL file
func (c *Collector) LogInteraction(sample domain.TrainingSample) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a daily log file: training_data_2025-12-24.jsonl
	filename := fmt.Sprintf("training_data_%s.jsonl", time.Now().Format("2006-01-02"))
	filePath := filepath.Join(c.DataDir, filename)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(sample)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}
	if _, err := f.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
