package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TaskRunRecord struct {
	TaskName   string `json:"task_name"`
	Action     string `json:"action"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	DurationMs int64  `json:"duration_ms"`
	Status     string `json:"status"`
	ExitCode   int    `json:"exit_code"`
	Message    string `json:"message,omitempty"`
}

func WriteTaskRun(path string, record TaskRunRecord) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create task run dir: %w", err)
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal task run: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write task run: %w", err)
	}
	return nil
}
