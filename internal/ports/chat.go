package ports

import (
	"context"
	"time"
)

// ChatService defines the unified communication interface for the Super Node.
// It bridges Matrix (Text/Channels) and LiveKit (Audio/Video).
type ChatService interface {
	// Matrix / Channel Operations
	CreateChannel(ctx context.Context, name, alias string, isPublic bool) (string, error)
	JoinChannel(ctx context.Context, channelID, userID string) error
	SendMessage(ctx context.Context, channelID, message string) error
	
	// LiveKit / Real-time Operations
	CreateRoom(ctx context.Context, roomName string) (string, error)
	GenerateToken(ctx context.Context, roomName, participantName string) (string, error)
	
	// Unified Operations (Link Matrix room to LiveKit call)
	StartCallInChannel(ctx context.Context, channelID string) (string, error)
}

// ChatMessage represents a standard message format across platforms
type ChatMessage struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Platform  string    `json:"platform"` // "matrix", "telegram", "whatsapp"
}
