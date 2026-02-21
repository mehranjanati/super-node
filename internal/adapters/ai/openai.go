package ai

import (
	"context"
	"errors"
	"io"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIClient wraps the official go-openai client
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient creates a new client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	if apiKey == "" {
		log.Println("Warning: OpenAI API Key is empty")
	}
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  openai.GPT4TurboPreview,
	}
}

// ChatRequest represents the incoming chat messages
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamChat sends a chat request to OpenAI and returns a channel of chunks
func (c *OpenAIClient) StreamChat(ctx context.Context, messages []ChatMessage) (<-chan string, <-chan error) {
	streamChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(streamChan)
		defer close(errChan)

		// Convert internal messages to OpenAI messages
		reqMessages := make([]openai.ChatCompletionMessage, len(messages))
		for i, m := range messages {
			reqMessages[i] = openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			}
		}

		// Add system prompt if not present (optional)
		// reqMessages = append([]openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: "You are VoltAgent..."}}, reqMessages...)

		req := openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: reqMessages,
			Stream:   true,
		}

		stream, err := c.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				errChan <- err
				return
			}

			if len(response.Choices) > 0 {
				chunk := response.Choices[0].Delta.Content
				if chunk != "" {
					streamChan <- chunk
				}
			}
		}
	}()

	return streamChan, errChan
}
