package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Record struct {
	JobName    string `json:"job_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	DurationMs int64  `json:"duration_ms"`
	ExitCode   int    `json:"exit_code"`
	StdoutPath string `json:"stdout_path"`
	StderrPath string `json:"stderr_path"`
	OS         string `json:"os"`
	Version    string `json:"binary_version"`
}

func WriteRecord(path string, record Record) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal record: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write record: %w", err)
	}
	return nil
}

func ReadRecord(path string) (Record, error) {
	var record Record
	data, err := os.ReadFile(path)
	if err != nil {
		return record, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return record, fmt.Errorf("parse record: %w", err)
	}
	return record, nil
}
