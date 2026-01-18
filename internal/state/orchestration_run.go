package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AgentRunRecord struct {
	Name       string `json:"name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	DurationMs int64  `json:"duration_ms"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type OrchestrationRunRecord struct {
	Orchestration string           `json:"orchestration"`
	Agents        []AgentRunRecord `json:"agents"`
	StartTime     string           `json:"start_time"`
	EndTime       string           `json:"end_time"`
	DurationMs    int64            `json:"duration_ms"`
	Status        string           `json:"status"`
}

func WriteOrchestrationRun(path string, record OrchestrationRunRecord) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create orchestration run dir: %w", err)
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal orchestration run: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write orchestration run: %w", err)
	}
	return nil
}

func ReadOrchestrationRun(path string) (OrchestrationRunRecord, error) {
	var record OrchestrationRunRecord
	data, err := os.ReadFile(path)
	if err != nil {
		return record, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return record, fmt.Errorf("parse orchestration run: %w", err)
	}
	return record, nil
}
