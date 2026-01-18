package taskflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/state"
)

func PlanCreate(name string, cfg config.Config, logger *logging.Logger, stateDir string, w io.Writer) (int, error) {
	if w == nil {
		w = io.Discard
	}
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := ValidateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	status := taskCfg.Status
	if status == "" {
		status = "planned"
	}

	now := time.Now().UTC()
	taskRecord := state.TaskRecord{
		Name:        name,
		Description: taskCfg.Description,
		Repo:        taskCfg.Repo,
		Status:      status,
		LastRun:     now.Format(time.RFC3339),
		Outputs:     taskCfg.Outputs,
		Documents:   taskCfg.Documents,
	}

	taskPath := filepath.Join(stateDir, "tasks", name+".json")
	if err := state.WriteTask(taskPath, taskRecord); err != nil {
		return 2, err
	}

	if err := WriteTaskRun(stateDir, name, "plan.create", now, now, status, 0, "created"); err != nil {
		logger.Error("failed to write plan run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(w, "task=%s status=%s\n", name, status)
	return 0, nil
}

func BuildRun(name string, cfg config.Config, logger *logging.Logger, stateDir string, w io.Writer) (int, error) {
	if w == nil {
		w = io.Discard
	}
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := ValidateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	start := time.Now().UTC()
	if err := UpdateTaskState(stateDir, name, taskCfg, "in_progress", start); err != nil {
		return 2, err
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, taskCfg.Command[0], taskCfg.Command[1:]...)
	cmd.Dir = taskCfg.WorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	logger.Info("task build starting", "task", name, "command", strings.Join(taskCfg.Command, " "))
	execErr := cmd.Run()
	end := time.Now().UTC()

	exitCode := exitCodeFromError(execErr)
	status := "done"
	message := "completed"
	if execErr != nil {
		status = "in_progress"
		message = execErr.Error()
		logger.Error("task build failed", "task", name, "error", execErr)
	}

	if err := UpdateTaskState(stateDir, name, taskCfg, status, end); err != nil {
		return 2, err
	}
	if err := WriteTaskRun(stateDir, name, "build.run", start, end, status, exitCode, message); err != nil {
		logger.Error("failed to write build run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(w, "task=%s exit=%d status=%s\n", name, exitCode, status)
	if execErr != nil {
		return 2, execErr
	}
	return 0, nil
}

func DocGenerate(name string, cfg config.Config, logger *logging.Logger, stateDir string, w io.Writer) (int, error) {
	if w == nil {
		w = io.Discard
	}
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := ValidateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	start := time.Now().UTC()
	status := "done"
	baseDir := taskCfg.WorkingDir

	if err := UpdateTaskState(stateDir, name, taskCfg, status, start); err != nil {
		return 2, err
	}
	if err := appendTaskSummary(baseDir, name, taskCfg, status); err != nil {
		logger.Error("failed to update README", "task", name, "error", err)
		return 2, err
	}

	docPath := filepath.Join(baseDir, "docs", "tasks", name+".md")
	if err := writeTaskDoc(docPath, name, taskCfg, status); err != nil {
		logger.Error("failed to write task doc", "task", name, "error", err)
		return 2, err
	}

	end := time.Now().UTC()
	if err := WriteTaskRun(stateDir, name, "doc.generate", start, end, status, 0, "documented"); err != nil {
		logger.Error("failed to write doc run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(w, "task=%s status=%s doc=%s\n", name, status, docPath)
	return 0, nil
}

func ResolveTaskStatus(stateDir string, name string, taskCfg config.TaskConfig) string {
	path := filepath.Join(stateDir, "tasks", name+".json")
	if record, err := state.ReadTask(path); err == nil && record.Status != "" {
		return record.Status
	}
	if taskCfg.Status != "" {
		return taskCfg.Status
	}
	return "planned"
}

func UpdateTaskState(stateDir string, name string, taskCfg config.TaskConfig, status string, lastRun time.Time) error {
	record := state.TaskRecord{
		Name:        name,
		Description: taskCfg.Description,
		Repo:        taskCfg.Repo,
		Status:      status,
		LastRun:     lastRun.Format(time.RFC3339),
		Outputs:     taskCfg.Outputs,
		Documents:   taskCfg.Documents,
	}
	path := filepath.Join(stateDir, "tasks", name+".json")
	return state.WriteTask(path, record)
}

func WriteTaskRun(stateDir string, taskName string, action string, start time.Time, end time.Time, status string, exitCode int, message string) error {
	timestamp := start.Format("20060102T150405Z")
	runPath := filepath.Join(stateDir, "runs", taskName, timestamp+".json")
	record := state.TaskRunRecord{
		TaskName:   taskName,
		Action:     action,
		StartTime:  start.Format(time.RFC3339),
		EndTime:    end.Format(time.RFC3339),
		DurationMs: end.Sub(start).Milliseconds(),
		Status:     status,
		ExitCode:   exitCode,
		Message:    message,
	}
	return state.WriteTaskRun(runPath, record)
}

func ValidateTaskConfig(name string, task config.TaskConfig) error {
	if name == "" {
		return errors.New("task name is required")
	}
	switch task.Repo {
	case "orchastration", "external":
	default:
		return fmt.Errorf("task %s has invalid repo: %s", name, task.Repo)
	}
	if task.WorkingDir == "" {
		return fmt.Errorf("task %s has empty working_dir", name)
	}
	if !filepath.IsAbs(task.WorkingDir) {
		return fmt.Errorf("task %s working_dir must be absolute", name)
	}
	if len(task.Command) == 0 {
		return fmt.Errorf("task %s has empty command", name)
	}
	return nil
}

func appendTaskSummary(baseDir string, name string, taskCfg config.TaskConfig, status string) error {
	section := buildTaskSummary(name, taskCfg, status)
	readmePath := filepath.Join(baseDir, "README.md")
	file, err := os.OpenFile(readmePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open README: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintln(file, "\n"+section); err != nil {
		return fmt.Errorf("append README: %w", err)
	}
	return nil
}

func writeTaskDoc(path string, name string, taskCfg config.TaskConfig, status string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create docs dir: %w", err)
	}
	content := buildTaskDoc(name, taskCfg, status)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write task doc: %w", err)
	}
	return nil
}

func buildTaskSummary(name string, taskCfg config.TaskConfig, status string) string {
	command := strings.Join(taskCfg.Command, " ")
	outputs := "(none)"
	if len(taskCfg.Outputs) > 0 {
		outputs = strings.Join(taskCfg.Outputs, ", ")
	}
	return fmt.Sprintf("## Task Summary: %s\n\n- Purpose: %s\n- Command: %s\n- Outputs: %s\n- Status: %s\n", name, taskCfg.Description, command, outputs, status)
}

func buildTaskDoc(name string, taskCfg config.TaskConfig, status string) string {
	command := strings.Join(taskCfg.Command, " ")
	outputs := "(none)"
	if len(taskCfg.Outputs) > 0 {
		outputs = strings.Join(taskCfg.Outputs, ", ")
	}
	documents := "(none)"
	if len(taskCfg.Documents) > 0 {
		documents = strings.Join(taskCfg.Documents, ", ")
	}
	return fmt.Sprintf("# Task: %s\n\n## Purpose\n%s\n\n## Commands Run\n%s\n\n## Outputs Produced\n%s\n\n## Documents\n%s\n\n## Status\n%s\n", name, taskCfg.Description, command, outputs, documents, status)
}

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}
