package social

import (
	"context"
	"testing"

	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSocialRepository
type MockSocialRepository struct {
	mock.Mock
}

// Ensure MockSocialRepository implements ports.SocialRepository
var _ ports.SocialRepository = (*MockSocialRepository)(nil)

func (m *MockSocialRepository) SavePost(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockSocialRepository) GetPosts(ctx context.Context, filter domain.FeedFilter) ([]*domain.Post, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func (m *MockSocialRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockSocialRepository) AddLike(ctx context.Context, postID, userID string) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

func (m *MockSocialRepository) SaveComment(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockSocialRepository) GetComments(ctx context.Context, postID string) ([]*domain.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*domain.Comment), args.Error(1)
}

// MockEventProducer
type MockEventProducer struct {
	mock.Mock
}

func (m *MockEventProducer) Produce(ctx context.Context, topic []byte, payload []byte) error {
	args := m.Called(ctx, topic, payload)
	return args.Error(0)
}

func (m *MockEventProducer) Close() error {
	return nil
}

// MockWasmRunner
type MockWasmRunner struct {
	mock.Mock
}

func (m *MockWasmRunner) RunModule(ctx context.Context, moduleID string, functionName string, params []uint64) ([]uint64, error) {
	args := m.Called(ctx, moduleID, functionName, params)
	return args.Get(0).([]uint64), args.Error(1)
}

func TestCreatePost(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	mockProd := new(MockEventProducer)
	mockWasm := new(MockWasmRunner)

	service := NewSocialService(mockRepo, mockProd, mockWasm)

	ctx := context.Background()
	authorID := "user-1"
	content := "Hello World"
	mediaURLs := []string{}

	// Expectations
	mockRepo.On("SavePost", ctx, mock.AnythingOfType("*domain.Post")).Return(nil)
	mockProd.On("Produce", ctx, []byte("social-feed"), mock.Anything).Return(nil)

	post, err := service.CreatePost(ctx, authorID, content, mediaURLs)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, content, post.Content)
	assert.Equal(t, authorID, post.AuthorID)

	mockRepo.AssertExpectations(t)
	mockProd.AssertExpectations(t)
}

func TestGetFeed(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	mockProd := new(MockEventProducer)
	mockWasm := new(MockWasmRunner)

	service := NewSocialService(mockRepo, mockProd, mockWasm)

	ctx := context.Background()
	filter := domain.FeedFilter{Limit: 10}
	expectedPosts := []*domain.Post{
		{ID: "post-1", Content: "Test 1"},
		{ID: "post-2", Content: "Test 2"},
	}

	mockRepo.On("GetPosts", ctx, filter).Return(expectedPosts, nil)

	posts, err := service.GetFeed(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "post-1", posts[0].ID)

	mockRepo.AssertExpectations(t)
}
