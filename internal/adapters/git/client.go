package git

import (
	"context"
	"fmt"
	"log"
	"nexus-super-node-v3/internal/ports"
	"time"
)

type MockGitClient struct {
	BaseURL string
	Token   string
}

func NewMockGitClient(token string) ports.GitProvider {
	return &MockGitClient{
		BaseURL: "https://github.com/nexus-user",
		Token:   token,
	}
}

func (c *MockGitClient) InitializeRepo(ctx context.Context, name string) (string, error) {
	log.Printf("[GitOps] Initializing repository: %s/%s", c.BaseURL, name)
	// Simulate network delay
	time.Sleep(500 * time.Millisecond)
	return fmt.Sprintf("%s/%s", c.BaseURL, name), nil
}

func (c *MockGitClient) CommitAndPush(ctx context.Context, repoPath string, files map[string]string, message string, branch string) error {
	log.Printf("[GitOps] Committing %d files to %s (branch: %s)", len(files), repoPath, branch)
	for path := range files {
		log.Printf("  + %s", path)
	}
	log.Printf("  Message: %s", message)
	time.Sleep(1 * time.Second)
	return nil
}

func (c *MockGitClient) CreatePullRequest(ctx context.Context, repoName string, title string, headBranch string, baseBranch string) (string, error) {
	log.Printf("[GitOps] Creating PR for %s: '%s' (%s -> %s)", repoName, title, headBranch, baseBranch)
	return fmt.Sprintf("https://github.com/nexus-user/%s/pull/1", repoName), nil
}

func (c *MockGitClient) GetCloneURL(repoName string) string {
	return fmt.Sprintf("%s/%s.git", c.BaseURL, repoName)
}

func (c *MockGitClient) GetLatestWorkflowRun(ctx context.Context, repoName string, branch string) (string, string, error) {
	// Simulate "completed" and "success" after some time
	// In a real mock, we might check how long ago the repo was initialized
	log.Printf("[GitOps] Checking workflow status for %s on branch %s", repoName, branch)
	return "completed", "success", nil
}

func (c *MockGitClient) DownloadArtifact(ctx context.Context, repoName string, artifactName string, destPath string) error {
	log.Printf("[GitOps] Downloading artifact '%s' from %s to %s", artifactName, repoName, destPath)
	// Simulate download delay
	time.Sleep(2 * time.Second)
	// Create a dummy file
	// In real implementation, use os.WriteFile
	return nil
}
