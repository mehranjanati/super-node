package benthos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
)

// Client interacts with the Redpanda Connect (Benthos) Streams API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Benthos client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// StreamConfig represents a minimal Benthos stream configuration
type StreamConfig struct {
	Input    map[string]interface{} `yaml:"input"`
	Pipeline PipelineConfig         `yaml:"pipeline"`
	Output   map[string]interface{} `yaml:"output"`
}

type PipelineConfig struct {
	Processors []map[string]interface{} `yaml:"processors"`
}

// DeployStream creates or updates a stream
func (c *Client) DeployStream(ctx context.Context, id string, config StreamConfig) error {
	body, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	url := fmt.Sprintf("%s/streams/%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/yaml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to deploy stream (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// DeleteStream removes a stream
func (c *Client) DeleteStream(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/streams/%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return fmt.Errorf("failed to delete stream (status %d)", resp.StatusCode)
	}

	return nil
}
