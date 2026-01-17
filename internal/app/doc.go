package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
)

func runDoc(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("doc requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "generate":
		return docGenerate(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown doc subcommand: %s", sub)
	}
}

func docGenerate(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("doc generate requires a task name")
	}

	name := args[0]
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := validateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	start := time.Now().UTC()
	status := "done"
	if err := updateTaskState(stateDir, name, taskCfg, status, start); err != nil {
		return 2, err
	}

	if err := appendTaskSummary(name, taskCfg, status); err != nil {
		logger.Error("failed to update README", "task", name, "error", err)
		return 2, err
	}

	docPath := filepath.Join("docs", "tasks", name+".md")
	if err := writeTaskDoc(docPath, name, taskCfg, status); err != nil {
		logger.Error("failed to write task doc", "task", name, "error", err)
		return 2, err
	}

	end := time.Now().UTC()
	if err := writeTaskRun(stateDir, name, "doc.generate", start, end, status, 0, "documented"); err != nil {
		logger.Error("failed to write doc run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(os.Stdout, "task=%s status=%s doc=%s\n", name, status, docPath)
	return 0, nil
}

func appendTaskSummary(name string, taskCfg config.TaskConfig, status string) error {
	section := buildTaskSummary(name, taskCfg, status)
	file, err := os.OpenFile("README.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open README: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("\n" + section + "\n"); err != nil {
		return fmt.Errorf("write README: %w", err)
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
