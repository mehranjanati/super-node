package chat

import (
	"context"
	"fmt"
	"time"

	"nexus-super-node-v3/internal/adapters/openclaw"
	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/ports"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// UnifiedChatService implements ports.ChatService
// It manages both Matrix (via OpenClaw) and LiveKit.
type UnifiedChatService struct {
	config       *config.Config
	lkRoomClient *lksdk.RoomServiceClient
	claw         *openclaw.Client
}

// Ensure implementation
var _ ports.ChatService = (*UnifiedChatService)(nil)

// NewUnifiedChatService creates a new instance
func NewUnifiedChatService(cfg *config.Config, claw *openclaw.Client) *UnifiedChatService {
	// Initialize LiveKit Client
	roomClient := lksdk.NewRoomServiceClient(
		cfg.LiveKit.APIURL,
		cfg.LiveKit.APIKey,
		cfg.LiveKit.APISecret,
	)

	svc := &UnifiedChatService{
		config:       cfg,
		lkRoomClient: roomClient,
		claw:         claw,
	}

	return svc
}

// --- LiveKit Implementation ---

func (s *UnifiedChatService) CreateRoom(ctx context.Context, roomName string) (string, error) {
	// LiveKit creates rooms on demand usually, but we can explicitly create one
	_, err := s.lkRoomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name: roomName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create LiveKit room: %w", err)
	}
	return roomName, nil
}

func (s *UnifiedChatService) GenerateToken(ctx context.Context, roomName, participantName string) (string, error) {
	at := auth.NewAccessToken(s.config.LiveKit.APIKey, s.config.LiveKit.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant).SetIdentity(participantName).SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

// --- Matrix Implementation (Delegated to OpenClaw) ---

// CreateChannel creates a new Matrix room.
// Note: OpenClaw currently does not expose an API to create rooms programmatically via the Gateway.
// Users should create rooms via their Matrix client (Element/OpenClaw UI) and invite the bot.
func (s *UnifiedChatService) CreateChannel(ctx context.Context, name, alias string, isPublic bool) (string, error) {
	return "", fmt.Errorf("creating channels programmatically is disabled. Please use OpenClaw/Element to create rooms and invite the bot")
}

// JoinChannel invites a user to a channel.
// Note: This relies on the bot being an admin in the room.
func (s *UnifiedChatService) JoinChannel(ctx context.Context, channelID, userID string) error {
	// If OpenClaw supports invites via API in the future, we can add it here.
	// For now, we assume the user joins via the room link or invite.
	return fmt.Errorf("programmatic invites are disabled. Please use OpenClaw/Element to invite users")
}

func (s *UnifiedChatService) SendMessage(ctx context.Context, channelID, message string) error {
	if s.claw == nil {
		return fmt.Errorf("OpenClaw client is not initialized")
	}

	// Use "matrix" channel explicitly as per user request
	return s.claw.SendMessage(channelID, message, "matrix")
}

// --- Unified Logic ---

// StartCallInChannel creates a LiveKit room for a Matrix channel and posts the link
func (s *UnifiedChatService) StartCallInChannel(ctx context.Context, channelID string) (string, error) {
	// 1. Create a LiveKit room with the same ID/Name as the channel
	roomName := fmt.Sprintf("call-%s", channelID)
	_, err := s.CreateRoom(ctx, roomName)
	if err != nil {
		return "", err
	}

	// 2. Generate a "Join Call" URL (In a real app, this would be a frontend URL)
	callURL := fmt.Sprintf("https://supernode.local/call/%s", roomName)

	// 3. Post the link to the Matrix channel
	msg := fmt.Sprintf("📞 A call has started! Join here: %s", callURL)
	if err := s.SendMessage(ctx, channelID, msg); err != nil {
		return "", err
	}

	return callURL, nil
}
