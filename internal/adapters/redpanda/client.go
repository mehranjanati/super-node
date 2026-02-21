package redpanda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/core/domain"

	"github.com/twmb/franz-go/pkg/kgo"
)

const redpandaMCPURL = "https://docs.redpanda.com/mcp"

// Client is a client for the Redpanda system (both MCP and Kafka).
type Client struct {
	url         string
	kafkaClient *kgo.Client
	config      *config.Config
}

// NewClient creates a new Redpanda client.
func NewClient(cfg *config.Config) (*Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Redpanda.Brokers...),
		kgo.ConsumerGroup(cfg.Redpanda.GroupID),
		kgo.ConsumeTopics(cfg.Redpanda.Topic),
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		// For development, if Redpanda isn't up, we might not want to crash immediately,
		// but for a real system we should probably error out or retry.
		// For now, let's log and return error.
		return nil, fmt.Errorf("failed to create redpanda client: %w", err)
	}

	return &Client{
		url:         redpandaMCPURL,
		kafkaClient: cl,
		config:      cfg,
	}, nil
}

// Produce sends a message to Redpanda.
func (c *Client) Produce(ctx context.Context, key, value []byte) error {
	if c.kafkaClient == nil {
		return fmt.Errorf("redpanda client not initialized")
	}

	record := &kgo.Record{
		Topic: c.config.Redpanda.Topic,
		Key:   key,
		Value: value,
	}

	return c.kafkaClient.ProduceSync(ctx, record).FirstErr()
}

// ConsumeLoop starts a consumer loop with the given handler.
// It supports high concurrency by processing messages in parallel if needed.
func (c *Client) ConsumeLoop(ctx context.Context, handler func(ctx context.Context, record *kgo.Record)) {
	if c.kafkaClient == nil {
		log.Println("Redpanda client not initialized, skipping consume loop")
		return
	}

	concurrency := c.config.Agents.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	// Semaphore to limit concurrency
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	log.Printf("Starting Redpanda consumer loop with concurrency: %d", concurrency)

	for {
		fetches := c.kafkaClient.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("Poll errors: %v", errs)
			// Don't exit loop on transient errors
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			sem <- struct{}{} // Acquire token
			wg.Add(1)

			go func(rec *kgo.Record) {
				defer func() {
					<-sem // Release token
					wg.Done()
				}()

				// Process message
				// In a real app, you might want a timeout per message
				handlerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				handler(handlerCtx, rec)
			}(record)
		}

		// Optional: Wait for batch to finish if strict ordering isn't required but
		// you want to commit offsets safely.
		// With Franz Go, auto-commit is enabled by default.
		// For high throughput, we usually let it auto-commit.
	}
}

// Close closes the underlying Kafka client.
func (c *Client) Close() {
	if c.kafkaClient != nil {
		c.kafkaClient.Close()
	}
}

// --- MCP Methods ---

func (c *Client) GetToolBelt(ctx context.Context) ([]*domain.MCPTool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting toolbelt: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback for demo if URL is unreachable
		return []*domain.MCPTool{}, nil
	}

	var toolbelt []*domain.MCPTool
	if err := json.NewDecoder(resp.Body).Decode(&toolbelt); err != nil {
		return nil, fmt.Errorf("error decoding toolbelt: %w", err)
	}

	return toolbelt, nil
}

func (c *Client) RouteToolCall(ctx context.Context, toolID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Implementation...
	return nil, nil
}
