package ports

import "context"

// GitProvider defines operations for version control integration
type GitProvider interface {
	// InitializeRepo creates a new repository or initializes one locally
	InitializeRepo(ctx context.Context, name string) (string, error)

	// CommitAndPush stages files, commits them, and pushes to a remote branch
	CommitAndPush(ctx context.Context, repoPath string, files map[string]string, message string, branch string) error

	// CreatePullRequest creates a PR for review
	CreatePullRequest(ctx context.Context, repoName string, title string, headBranch string, baseBranch string) (string, error)

	// GetCloneURL returns the HTTP clone URL for a repository
	GetCloneURL(repoName string) string

	// GetLatestWorkflowRun returns the status of the latest workflow run for a branch
	GetLatestWorkflowRun(ctx context.Context, repoName string, branch string) (string, string, error) // status, conclusion, error

	// DownloadArtifact downloads the build artifact from the latest successful run
	DownloadArtifact(ctx context.Context, repoName string, artifactName string, destPath string) error
}
