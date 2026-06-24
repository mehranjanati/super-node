package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"nexus-super-node-v3/internal/core/domain"
	workflowpkg "nexus-super-node-v3/internal/workflow"
)

const (
	workflowExecutionIndexKey  = "workflow_executions:index"
	workflowExecutionKeyPrefix = "workflow_execution:"
)

type workflowExecutionIndex struct {
	IDs []string `json:"ids"`
}

type workflowLogEntry struct {
	WorkflowID  string    `json:"workflow_id"`
	Time        time.Time `json:"time"`
	Level       string    `json:"level"`
	Service     string    `json:"service"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CurrentStep string    `json:"current_step"`
}

func workflowExecutionKey(workflowID string) string {
	return workflowExecutionKeyPrefix + workflowID
}

func (g *EchoGateway) listWorkflowExecutionRecords(ctx context.Context) ([]domain.WorkflowExecutionRecord, error) {
	if g.appDataRepo == nil {
		return nil, nil
	}

	index, err := g.loadWorkflowExecutionIndex(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]domain.WorkflowExecutionRecord, 0, len(index.IDs))
	for _, workflowID := range index.IDs {
		record, err := g.getWorkflowExecutionRecord(ctx, workflowID)
		if err != nil {
			return nil, err
		}
		if record != nil {
			records = append(records, *record)
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].UpdatedAt.After(records[j].UpdatedAt)
	})

	return records, nil
}

func (g *EchoGateway) getWorkflowExecutionRecord(ctx context.Context, workflowID string) (*domain.WorkflowExecutionRecord, error) {
	if g.appDataRepo == nil {
		return nil, nil
	}

	item, err := g.appDataRepo.GetAppData(ctx, workflowExecutionKey(workflowID))
	if err != nil || item == nil {
		return nil, err
	}

	var record domain.WorkflowExecutionRecord
	if err := json.Unmarshal(item.Data, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (g *EchoGateway) saveWorkflowExecutionRecord(ctx context.Context, record domain.WorkflowExecutionRecord) error {
	if g.appDataRepo == nil {
		return nil
	}

	now := time.Now().UTC()
	if record.CreatedAt.IsZero() {
		record.CreatedAt = now
	}
	record.UpdatedAt = now
	record.Logs = compactLogs(record.Logs)

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if appErr := g.appDataRepo.UpsertAppData(ctx, workflowExecutionKey(record.WorkflowID), payload); appErr != nil {
		return appErr
	}

	index, err := g.loadWorkflowExecutionIndex(ctx)
	if err != nil {
		return err
	}

	exists := false
	for _, existingID := range index.IDs {
		if existingID == record.WorkflowID {
			exists = true
			break
		}
	}
	if !exists {
		index.IDs = append([]string{record.WorkflowID}, index.IDs...)
	}

	indexPayload, err := json.Marshal(index)
	if err != nil {
		return err
	}

	return g.appDataRepo.UpsertAppData(ctx, workflowExecutionIndexKey, indexPayload)
}

func (g *EchoGateway) loadWorkflowExecutionIndex(ctx context.Context) (*workflowExecutionIndex, error) {
	if g.appDataRepo == nil {
		return &workflowExecutionIndex{}, nil
	}

	item, err := g.appDataRepo.GetAppData(ctx, workflowExecutionIndexKey)
	if err != nil || item == nil {
		return &workflowExecutionIndex{}, err
	}

	var index workflowExecutionIndex
	if err := json.Unmarshal(item.Data, &index); err != nil {
		return nil, err
	}

	return &index, nil
}

func (g *EchoGateway) persistWorkflowStart(ctx context.Context, workflowID, kind, name, currentStep string, logs []string, artifacts *domain.WorkflowExecutionArtifacts) error {
	existing, err := g.getWorkflowExecutionRecord(ctx, workflowID)
	if err != nil {
		return err
	}

	record := domain.WorkflowExecutionRecord{
		WorkflowID:  workflowID,
		Kind:        kind,
		Name:        name,
		Status:      "RUNNING",
		CurrentStep: currentStep,
		Logs:        logs,
		Artifacts:   artifacts,
	}
	if existing != nil {
		record.CreatedAt = existing.CreatedAt
		if record.Kind == "" {
			record.Kind = existing.Kind
		}
		if record.Name == "" {
			record.Name = existing.Name
		}
		if record.Artifacts == nil {
			record.Artifacts = existing.Artifacts
		}
	}

	return g.saveWorkflowExecutionRecord(ctx, record)
}

func (g *EchoGateway) persistWorkflowStatus(ctx context.Context, workflowID string, status *workflowpkg.DummyDeploymentStatus) (*domain.WorkflowExecutionRecord, error) {
	existing, err := g.getWorkflowExecutionRecord(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	record := domain.WorkflowExecutionRecord{
		WorkflowID:  workflowID,
		Kind:        "deployment",
		Name:        workflowID,
		Status:      normalizeWorkflowStatus(status.Status),
		CurrentStep: status.CurrentStep,
		Logs:        status.Logs,
		Error:       status.Error,
	}

	if existing != nil {
		record.CreatedAt = existing.CreatedAt
		if existing.Kind != "" {
			record.Kind = existing.Kind
		}
		if existing.Name != "" {
			record.Name = existing.Name
		}
		record.Artifacts = existing.Artifacts
	}

	if status.Artifacts != nil {
		record.Artifacts = &domain.WorkflowExecutionArtifacts{
			ProjectName:    status.Artifacts.ProjectName,
			Template:       status.Artifacts.Template,
			PlanningSource: existingPlanningSource(existing),
			URL:            status.Artifacts.URL,
			Message:        status.Artifacts.Message,
		}
		if record.Name == workflowID {
			record.Name = firstNonEmpty(status.Artifacts.ProjectName, workflowID)
		}
	}

	if err := g.saveWorkflowExecutionRecord(ctx, record); err != nil {
		return nil, err
	}

	return g.getWorkflowExecutionRecord(ctx, workflowID)
}

func (g *EchoGateway) listWorkflowLogEntries(ctx context.Context) ([]workflowLogEntry, error) {
	records, err := g.listWorkflowExecutionRecords(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]workflowLogEntry, 0)
	for _, record := range records {
		when := record.UpdatedAt
		if when.IsZero() {
			when = record.CreatedAt
		}
		if when.IsZero() {
			when = time.Now().UTC()
		}

		for _, message := range record.Logs {
			entries = append(entries, workflowLogEntry{
				WorkflowID:  record.WorkflowID,
				Time:        when,
				Level:       logLevelForRecord(record),
				Service:     serviceForKind(record.Kind),
				Message:     message,
				Status:      record.Status,
				CurrentStep: record.CurrentStep,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.After(entries[j].Time)
	})

	return entries, nil
}

func compactLogs(logs []string) []string {
	compact := make([]string, 0, len(logs))
	for _, log := range logs {
		log = strings.TrimSpace(log)
		if log == "" {
			continue
		}
		compact = append(compact, log)
	}
	return compact
}

func normalizeWorkflowStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "COMPLETED", "DONE", "SUCCESS":
		return "COMPLETED"
	case "FAILED", "ERROR", "CANCELED", "TERMINATED", "TIMED_OUT":
		return "FAILED"
	default:
		return "RUNNING"
	}
}

func logLevelForRecord(record domain.WorkflowExecutionRecord) string {
	if normalizeWorkflowStatus(record.Status) == "FAILED" || record.Error != "" {
		return "ERROR"
	}
	if record.CurrentStep == "DONE" {
		return "INFO"
	}
	return "INFO"
}

func serviceForKind(kind string) string {
	switch kind {
	case "deployment":
		return "TEMPORAL_DEPLOY"
	case "crypto_analysis":
		return "TEMPORAL_CRYPTO"
	case "human_handoff":
		return "TEMPORAL_HANDOFF"
	default:
		return "TEMPORAL_WORKFLOW"
	}
}

func inferWorkflowName(toolID, workflowID string, args map[string]interface{}) string {
	switch toolID {
	case "system__deploy_website":
		return firstNonEmpty(stringArg(args, "project_name"), workflowID)
	case "system__crypto_analysis":
		return firstNonEmpty(fmt.Sprintf("crypto-%s", stringArg(args, "user_id")), workflowID)
	case "system__human_handoff":
		return firstNonEmpty(fmt.Sprintf("handoff-%s", stringArg(args, "room_id")), workflowID)
	default:
		return workflowID
	}
}

func inferWorkflowKind(toolID string) string {
	switch toolID {
	case "system__deploy_website":
		return "deployment"
	case "system__crypto_analysis":
		return "crypto_analysis"
	case "system__human_handoff":
		return "human_handoff"
	default:
		return "workflow"
	}
}

func inferWorkflowArtifacts(toolID string, args map[string]interface{}, result map[string]string) *domain.WorkflowExecutionArtifacts {
	switch toolID {
	case "system__deploy_website":
		return &domain.WorkflowExecutionArtifacts{
			ProjectName:    stringArg(args, "project_name"),
			PlanningSource: stringMapValue(result, "planning_source"),
			Message:        "Workflow queued via VoltAgent.",
		}
	default:
		return nil
	}
}

func existingPlanningSource(record *domain.WorkflowExecutionRecord) string {
	if record == nil || record.Artifacts == nil {
		return ""
	}
	return record.Artifacts.PlanningSource
}

func stringArg(args map[string]interface{}, key string) string {
	if value, ok := args[key].(string); ok {
		return value
	}
	return ""
}

func stringMapValue(values map[string]string, key string) string {
	if values == nil {
		return ""
	}
	return strings.TrimSpace(values[key])
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
