package openclaw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client handles communication with OpenClaw Gateway
type Client struct {
	BaseURL    string
	AuthSecret string
	Client     *http.Client
}

// NewClient creates a new OpenClaw client
func NewClient(baseURL, authSecret string) *Client {
	return &Client{
		BaseURL:    baseURL,
		AuthSecret: authSecret,
		Client:     &http.Client{},
	}
}

// SendMessageRequest represents the payload for sending a message
type SendMessageRequest struct {
	To      string `json:"to"`      // E.g., matrix user ID or phone number
	Message string `json:"message"` // The text content
	Channel string `json:"channel"` // Optional: specific channel (e.g., "matrix", "whatsapp")
}

// SendMessage sends a message to a user via OpenClaw
func (c *Client) SendMessage(to, message, channel string) error {
	payload := SendMessageRequest{
		To:      to,
		Message: message,
		Channel: channel,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/send", c.BaseURL) // Assuming OpenClaw has a simple send endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("openclaw returned status: %d", resp.StatusCode)
	}

	return nil
}
