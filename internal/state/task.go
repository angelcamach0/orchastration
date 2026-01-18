package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TaskRecord struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Repo        string   `json:"repo"`
	Status      string   `json:"status"`
	LastRun     string   `json:"last_run"`
	Outputs     []string `json:"outputs"`
	Documents   []string `json:"documents"`
}

func WriteTask(path string, task TaskRecord) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create task dir: %w", err)
	}

	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal task: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write task: %w", err)
	}
	return nil
}

func ReadTask(path string) (TaskRecord, error) {
	var task TaskRecord
	data, err := os.ReadFile(path)
	if err != nil {
		return task, err
	}
	if err := json.Unmarshal(data, &task); err != nil {
		return task, fmt.Errorf("parse task: %w", err)
	}
	return task, nil
}
