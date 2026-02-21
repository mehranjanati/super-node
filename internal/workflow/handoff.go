package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/livekit/protocol/auth"
)

// Activities
type HandoffActivities struct {
	LiveKitAPIKey    string
	LiveKitAPISecret string
}

func (a *HandoffActivities) LogToTerminal(ctx context.Context, message string) error {
	fmt.Println("[Orchestrator] " + message)
	return nil
}

type InitiateHandoffParams struct {
	RoomID       string
	TargetUserID string
}

func (a *HandoffActivities) SendHandoffEvent(ctx context.Context, params InitiateHandoffParams) error {
	// Mock Event
	fmt.Println("[Handoff Event] Triggering LiveKit room: " + params.RoomID)

	// Generate Token
	at := auth.NewAccessToken(a.LiveKitAPIKey, a.LiveKitAPISecret)
	canPublish := true
	canSubscribe := true
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         params.RoomID,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	at.AddGrant(grant).SetIdentity(params.TargetUserID)

	token, err := at.ToJWT()
	if err != nil {
		return err
	}

	content := map[string]interface{}{
		"livekit_room": params.RoomID,
		"target_user":  params.TargetUserID,
		"timestamp":    time.Now().UnixMilli(),
		"token":        token,
		"msgtype":      "m.text", // Fallback
		"body":         "Incoming Call",
	}

	fmt.Printf("[Handoff] Notification content: %+v\n", content)
	return nil
}

func (a *HandoffActivities) TriggerTelegram(ctx context.Context, userId string) error {
	// Mock
	fmt.Printf("[Telegram] Bridging call to %s...\n", userId)
	return nil
}

func (a *HandoffActivities) TriggerSIP(ctx context.Context, userId string) error {
	// Mock
	fmt.Printf("[SIP] Dialing mobile number for %s...\n", userId)
	return nil
}
