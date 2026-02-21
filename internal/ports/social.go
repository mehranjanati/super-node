package ports

import (
	"context"

	"nexus-super-node-v3/internal/core/domain"
)

// SocialService defines the business logic for the DPIN Social Feed
type SocialService interface {
	CreatePost(ctx context.Context, authorID, content string, mediaURLs []string) (*domain.Post, error)
	GetFeed(ctx context.Context, filter domain.FeedFilter) ([]*domain.Post, error)
	LikePost(ctx context.Context, postID, userID string) error
	AddComment(ctx context.Context, postID, userID, content string) (*domain.Comment, error)
}

// SocialRepository defines persistence for social data
// Note: In a real DPIN, this might be a mix of Postgres + IPFS/Arweave
type SocialRepository interface {
	SavePost(ctx context.Context, post *domain.Post) error
	GetPosts(ctx context.Context, filter domain.FeedFilter) ([]*domain.Post, error)
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	AddLike(ctx context.Context, postID, userID string) error
	SaveComment(ctx context.Context, comment *domain.Comment) error
}
