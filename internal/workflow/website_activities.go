package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/benthos"
	"nexus-super-node-v3/internal/core/services/wazero"
	"nexus-super-node-v3/internal/ports"
	"os"
	"time"
)

type WebsiteDeploymentParams struct {
	ProjectName string `json:"project_name"`
	Prompt      string `json:"prompt"`
	Theme       string `json:"theme"`
	Framework   string `json:"framework"` // e.g., "svelte", "react"
}

type WebsiteDeploymentResult struct {
	URL       string `json:"url"`
	Status    string `json:"status"`
	ProjectID string `json:"project_id"`
	RepoURL   string `json:"repo_url"`
	PRURL     string `json:"pr_url"`
}

type WebsiteActivities struct {
	GitProvider   ports.GitProvider
	WazeroSvc     *wazero.WazeroService
	AIClient      *ai.OpenAIClient
	BenthosClient *benthos.Client
}

func (a *WebsiteActivities) GenerateUISchema(ctx context.Context, params WebsiteDeploymentParams) (string, error) {
	// In a real scenario, this would call an LLM with v0-like prompts
	return fmt.Sprintf(`{"components": [{"type": "hero", "title": "Welcome to %s", "theme": "%s"}]}`, params.ProjectName, params.Theme), nil
}

func (a *WebsiteActivities) GenerateSourceCode(ctx context.Context, schema string) (map[string]string, error) {
	if a.AIClient == nil {
		// Fallback for tests or if AI is disabled
		return GenerateProjectStructure("demo-site", "Basic Schema"), nil
	}

	// Use AI to generate code
	// In a real scenario, we would stream this or handle larger context
	// For now, we simulate a prompt to the AI
	_ = fmt.Sprintf("Generate a simple SvelteKit project structure based on this schema: %s. Return ONLY JSON with file paths as keys and content as values.", schema)

	// We use a simplified non-streaming call here for the activity
	// Ideally, we'd add a Non-Streaming method to OpenAIClient or consume the stream
	// For this example, we'll just mock the AI response based on the presence of the client
	// to avoid blocking on a stream implementation detail in this step.

	// REAL IMPLEMENTATION WOULD BE:
	// response, err := a.AIClient.ChatCompletion(ctx, prompt)

	// Extract project name from schema or use fallback
	var schemaMap map[string]interface{}
	projectName := "generated-site"
	if err := json.Unmarshal([]byte(schema), &schemaMap); err == nil {
		if name, ok := schemaMap["project_name"].(string); ok && name != "" {
			projectName = name
		} else if name, ok := schemaMap["name"].(string); ok && name != "" {
			projectName = name
		}
	}

	// Use the template generator to create a full stack structure (Frontend + WASM Backend)
	files := GenerateProjectStructure(projectName, schema)

	// If AI is available, we could enhance specific files here,
	// but the base structure is now robust and compilation-ready.
	if a.AIClient != nil {
		// Example: Enhance the Svelte component with the specific prompt details
		// enhancedSvelte := a.AIClient.EnhanceCode(files["frontend/src/routes/+page.svelte"], prompt)
		// files["frontend/src/routes/+page.svelte"] = enhancedSvelte
	}

	return files, nil
}

func (a *WebsiteActivities) PushToRepository(ctx context.Context, params WebsiteDeploymentParams, files map[string]string) (string, error) {
	if a.GitProvider == nil {
		return "", fmt.Errorf("git provider not configured")
	}

	repoName := fmt.Sprintf("nexus-app-%s", params.ProjectName)
	repoURL, err := a.GitProvider.InitializeRepo(ctx, repoName)
	if err != nil {
		return "", err
	}

	branch := "feature/init"
	err = a.GitProvider.CommitAndPush(ctx, repoURL, files, "Initial commit from Nexus Super Node", branch)
	if err != nil {
		return "", err
	}

	// Create PR for review
	prURL, err := a.GitProvider.CreatePullRequest(ctx, repoName, "feat: Initial App Generation", branch, "main")
	if err != nil {
		return repoURL, nil // Return repo URL even if PR fails
	}

	return prURL, nil
}

func (a *WebsiteActivities) BuildWebsiteBundle(ctx context.Context, projectName string) (string, error) {
	// Trigger build by checking CI status (assuming PushToRepository triggered it)
	// In a real scenario, we might trigger it manually via API if not automatic

	repoName := fmt.Sprintf("nexus-app-%s", projectName)
	branch := "feature/init"
	artifactName := "wasm-bundle"

	// Poll for completion
	maxRetries := 30
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		status, conclusion, err := a.GitProvider.GetLatestWorkflowRun(ctx, repoName, branch)
		if err != nil {
			// Log error and continue polling, maybe transient
			fmt.Printf("Error checking workflow status (attempt %d): %v\n", i+1, err)
			time.Sleep(retryInterval)
			continue
		}

		if status == "completed" {
			if conclusion == "success" {
				// Save to shared volume so wasm-processor can access it
				artifactPath := fmt.Sprintf("/data/wasm/bundle-%s.wasm", projectName)
				err = a.GitProvider.DownloadArtifact(ctx, repoName, artifactName, artifactPath)
				if err != nil {
					return "", fmt.Errorf("failed to download artifact: %w", err)
				}
				return artifactPath, nil
			} else {
				return "", fmt.Errorf("build workflow failed with conclusion: %s", conclusion)
			}
		}

		// If not completed, wait and retry
		fmt.Printf("Build in progress (status: %s)... waiting %v\n", status, retryInterval)
		time.Sleep(retryInterval)
	}

	return "", fmt.Errorf("build timed out after %d attempts", maxRetries)
}

func (a *WebsiteActivities) DeployToHosting(ctx context.Context, bundlePath string) (WebsiteDeploymentResult, error) {
	// Validate Wasm first
	wasmBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		return WebsiteDeploymentResult{Status: "failed"}, fmt.Errorf("failed to read wasm bundle: %w", err)
	}

	if a.WazeroSvc != nil {
		_, err := a.WazeroSvc.ExecuteModule(ctx, wasmBytes, "validate")
		if err != nil {
			fmt.Printf("Warning: validation failed: %v\n", err)
		}
	}

	projectID := fmt.Sprintf("site-%d", time.Now().Unix())

	// Deploy to Redpanda Connect (Benthos)
	// We use the projectID as the stream ID
	if a.BenthosClient != nil {
		err := a.DeployAgent(ctx, bundlePath, projectID)
		if err != nil {
			fmt.Printf("Warning: failed to deploy to Benthos: %v\n", err)
			// Don't fail the workflow, just log it for now
		}
	}

	return WebsiteDeploymentResult{
		URL:       fmt.Sprintf("http://%s.nexus.local", projectID),
		ProjectID: projectID,
	}, nil
}

// DeployAgent deploys a Wasm agent to the Redpanda Connect runtime
func (a *WebsiteActivities) DeployAgent(ctx context.Context, bundlePath string, agentID string) error {
	if a.BenthosClient == nil {
		return fmt.Errorf("benthos client not configured")
	}

	// Construct the config
	// Note: We point to /data/wasm because that's where the volume is mounted in the wasm-processor container too
	// (See docker-compose.yml: wasm_artifacts:/data/wasm)
	config := benthos.StreamConfig{
		Input: map[string]interface{}{
			"kafka": map[string]interface{}{
				"addresses":      []string{"redpanda:29092"},
				"topics":         []string{fmt.Sprintf("agent-%s-input", agentID)},
				"consumer_group": fmt.Sprintf("agent-group-%s", agentID),
			},
		},
		Pipeline: benthos.PipelineConfig{
			Processors: []map[string]interface{}{
				{
					"wasm": map[string]interface{}{
						"module_path": bundlePath, // Shared path
						"function":    "process",
					},
				},
			},
		},
		Output: map[string]interface{}{
			"kafka": map[string]interface{}{
				"addresses": []string{"redpanda:29092"},
				"topic":     fmt.Sprintf("agent-%s-output", agentID),
			},
		},
	}

	return a.BenthosClient.DeployStream(ctx, agentID, config)
}
