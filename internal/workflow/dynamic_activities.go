package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"nexus-super-node-v3/internal/core/services/mcp"
)

// DynamicActivities wraps typed activities to accept generic maps
type DynamicActivities struct {
	WebsiteActivities *WebsiteActivities
	CryptoActivities  *CryptoActivities
	HandoffActivities *HandoffActivities
}

func (d *DynamicActivities) LogThoughtActivityWrapper(ctx context.Context, args map[string]interface{}) error {
	roomID, _ := args["room_id"].(string)
	agentID, _ := args["agent_id"].(string)
	thought, _ := args["thought"].(string)
	return d.CryptoActivities.LogThoughtActivity(ctx, roomID, agentID, thought)
}

func (d *DynamicActivities) FetchMarketDataActivityWrapper(ctx context.Context, args map[string]interface{}) ([]map[string]interface{}, error) {
	return d.CryptoActivities.FetchMarketDataActivity(ctx)
}

func (d *DynamicActivities) AnalyzeMarketActivityWrapper(ctx context.Context, args map[string]interface{}) (mcp.CryptoAnalysisResult, error) {
	marketDataRaw, ok := args["market_data"].([]interface{})
	if !ok {
		return mcp.CryptoAnalysisResult{}, fmt.Errorf("missing or invalid 'market_data' argument")
	}

	marketData := make([]map[string]interface{}, len(marketDataRaw))
	for i, v := range marketDataRaw {
		marketData[i] = v.(map[string]interface{})
	}

	return d.CryptoActivities.AnalyzeMarketActivity(ctx, marketData)
}

func (d *DynamicActivities) NotifyUserActivityWrapper(ctx context.Context, args map[string]interface{}) error {
	userID, _ := args["user_id"].(string)
	analysisRaw, _ := json.Marshal(args["analysis"])
	var analysis mcp.CryptoAnalysisResult
	json.Unmarshal(analysisRaw, &analysis)

	return d.CryptoActivities.NotifyUserActivity(ctx, userID, analysis)
}

func (d *DynamicActivities) ExecuteTradeActivityWrapper(ctx context.Context, args map[string]interface{}) (string, error) {
	token, _ := args["token"].(string)
	amount, _ := args["amount"].(float64)
	return d.CryptoActivities.ExecuteTradeActivity(ctx, token, amount)
}

func (d *DynamicActivities) LogToTerminalWrapper(ctx context.Context, args map[string]interface{}) error {
	message, _ := args["message"].(string)
	return d.HandoffActivities.LogToTerminal(ctx, message)
}

func (d *DynamicActivities) SendHandoffEventWrapper(ctx context.Context, args map[string]interface{}) error {
	var params InitiateHandoffParams
	jsonRaw, _ := json.Marshal(args)
	json.Unmarshal(jsonRaw, &params)
	return d.HandoffActivities.SendHandoffEvent(ctx, params)
}

func (d *DynamicActivities) GenerateUISchemaWrapper(ctx context.Context, args map[string]interface{}) (string, error) {
	// Convert map to struct
	jsonBytes, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	var params WebsiteDeploymentParams
	if err := json.Unmarshal(jsonBytes, &params); err != nil {
		return "", fmt.Errorf("invalid args for GenerateUISchema: %w", err)
	}
	return d.WebsiteActivities.GenerateUISchema(ctx, params)
}

func (d *DynamicActivities) GenerateSourceCodeWrapper(ctx context.Context, args map[string]interface{}) (map[string]string, error) {
	// Simple argument mapping: assume "schema" key or raw string?
	// The workflow resolver passes a map.
	schema, ok := args["schema"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'schema' argument")
	}
	return d.WebsiteActivities.GenerateSourceCode(ctx, schema)
}

func (d *DynamicActivities) PushToRepositoryWrapper(ctx context.Context, args map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	var params WebsiteDeploymentParams
	if err := json.Unmarshal(jsonBytes, &params); err != nil {
		return "", fmt.Errorf("invalid args for PushToRepository: %w", err)
	}

	// Let's assume files are in "files" key
	var files map[string]string

	if f, ok := args["files"].(map[string]string); ok {
		files = f
	} else if fRaw, ok := args["files"].(map[string]interface{}); ok {
		files = make(map[string]string)
		for k, v := range fRaw {
			if s, ok := v.(string); ok {
				files[k] = s
			}
		}
	} else {
		return "", fmt.Errorf("missing or invalid 'files' argument: %T", args["files"])
	}

	return d.WebsiteActivities.PushToRepository(ctx, params, files)
}

func (d *DynamicActivities) BuildWebsiteBundleWrapper(ctx context.Context, args map[string]interface{}) (string, error) {
	projectName, ok := args["project_name"].(string)
	if !ok {
		return "", fmt.Errorf("missing 'project_name' argument")
	}
	return d.WebsiteActivities.BuildWebsiteBundle(ctx, projectName)
}

func (d *DynamicActivities) DeployToHostingWrapper(ctx context.Context, args map[string]interface{}) (WebsiteDeploymentResult, error) {
	bundlePath, ok := args["bundle_path"].(string)
	if !ok {
		return WebsiteDeploymentResult{Status: "failed"}, fmt.Errorf("missing 'bundle_path' argument")
	}
	return d.WebsiteActivities.DeployToHosting(ctx, bundlePath)
}
