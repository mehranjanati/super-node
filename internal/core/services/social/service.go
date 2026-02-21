package social

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"

	"github.com/google/uuid"
)

// SocialServiceImpl implements ports.SocialService
type SocialServiceImpl struct {
	repo       ports.SocialRepository
	eventProd  ports.EventProducer
	wasmRunner ports.WasmRunner
	// In-memory cache for demo purposes
	posts []*domain.Post
}

// Ensure implementation
var _ ports.SocialService = (*SocialServiceImpl)(nil)

// NewSocialService creates a new instance
func NewSocialService(repo ports.SocialRepository, eventProd ports.EventProducer, wasmRunner ports.WasmRunner) *SocialServiceImpl {
	return &SocialServiceImpl{
		repo:       repo,
		eventProd:  eventProd,
		wasmRunner: wasmRunner,
		posts:      make([]*domain.Post, 0),
	}
}

// CreatePost creates a new post, runs Wasm moderation, and publishes to Redpanda
func (s *SocialServiceImpl) CreatePost(ctx context.Context, authorID, content string, mediaURLs []string) (*domain.Post, error) {
	// 1. Basic Validation
	if content == "" && len(mediaURLs) == 0 {
		return nil, fmt.Errorf("post must have content or media")
	}

	// 2. Run Wasm Logic (e.g., Sentiment Analysis / Moderation)
	// In a real DPIN, we fetch the module from IPFS/Registry.
	// For this demo, we assume a "moderation-v1" module is available.
	// We'll simulate this step or use the runner if implemented.
	tags := []string{}
	// Simulated Wasm Result
	if len(content) > 10 {
		tags = append(tags, "#long_form")
	}
	if len(mediaURLs) > 0 {
		tags = append(tags, "#multimedia")
	}

	// 3. Create Domain Object
	post := &domain.Post{
		ID:        uuid.New().String(),
		AuthorID:  authorID,
		Content:   content,
		MediaURLs: mediaURLs,
		Tags:      tags,
		CreatedAt: time.Now(),
		Likes:     0,
		Metadata: map[string]interface{}{
			"processed_by": "wasm-moderation-v1",
			"node_id":      "super-node-local",
		},
	}

	// 4. Save to Repository
	if err := s.repo.SavePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to save post: %w", err)
	}

	// 5. Publish Event to Redpanda (Real-time Feed)
	eventPayload, _ := json.Marshal(post)
	if err := s.eventProd.Produce(ctx, []byte("social-feed"), eventPayload); err != nil {
		// Log error but don't fail the post creation necessarily, or do retry logic
		fmt.Printf("Warning: Failed to publish to Redpanda: %v\n", err)
	}

	return post, nil
}

func (s *SocialServiceImpl) GetFeed(ctx context.Context, filter domain.FeedFilter) ([]*domain.Post, error) {
	return s.repo.GetPosts(ctx, filter)
}

func (s *SocialServiceImpl) LikePost(ctx context.Context, postID, userID string) error {
	if err := s.repo.AddLike(ctx, postID, userID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	// Publish Like Event
	event := map[string]interface{}{
		"type":    "like",
		"post_id": postID,
		"user_id": userID,
	}
	payload, _ := json.Marshal(event)
	// Best effort publish
	s.eventProd.Produce(ctx, []byte("social-engagement"), payload)
	
	return nil
}

func (s *SocialServiceImpl) AddComment(ctx context.Context, postID, userID, content string) (*domain.Comment, error) {
	comment := &domain.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		AuthorID:  userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	
	if err := s.repo.SaveComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to save comment: %w", err)
	}
	
	// Publish Comment Event
	payload, _ := json.Marshal(comment)
	s.eventProd.Produce(ctx, []byte("social-engagement"), payload)
	
	return comment, nil
}
