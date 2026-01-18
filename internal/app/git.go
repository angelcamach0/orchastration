package app

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/state"
	"orchastration/internal/taskflow"
)

func runGit(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("git requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "issue":
		return gitIssue(args[1:], cfg, logger, stateDir)
	case "branch":
		return gitBranch(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown git subcommand: %s", sub)
	}
}

func gitIssue(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) < 2 || args[0] != "create" {
		return 2, errors.New("git issue create requires a task name")
	}

	name := args[1]
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := taskflow.ValidateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	title := fmt.Sprintf("Task: %s", name)
	body := taskCfg.Description
	if body == "" {
		body = "(no description provided)"
	}

	start := time.Now().UTC()
	cmd := exec.CommandContext(context.Background(), "gh", "issue", "create", "--title", title, "--body", body)
	cmd.Dir = taskCfg.WorkingDir
	output, err := cmd.CombinedOutput()
	end := time.Now().UTC()

	status := resolveTaskStatus(stateDir, name, taskCfg)
	if err := taskflow.UpdateTaskState(stateDir, name, taskCfg, status, end); err != nil {
		return 2, err
	}
	if runErr := taskflow.WriteTaskRun(stateDir, name, "git.issue.create", start, end, status, exitCodeFromError(err), strings.TrimSpace(string(output))); runErr != nil {
		logger.Error("failed to write git issue run", "task", name, "error", runErr)
	}

	if err != nil {
		return 2, fmt.Errorf("git issue create failed: %w", err)
	}
	return 0, nil
}

func gitBranch(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) < 2 || args[0] != "create" {
		return 2, errors.New("git branch create requires a task name")
	}

	name := args[1]
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := taskflow.ValidateTaskConfig(name, taskCfg); err != nil {
		return 2, err
	}

	branchName := buildBranchName(name, taskCfg.Repo)
	start := time.Now().UTC()
	cmd := exec.CommandContext(context.Background(), "git", "-C", taskCfg.WorkingDir, "checkout", "-b", branchName)
	output, err := cmd.CombinedOutput()
	end := time.Now().UTC()

	status := resolveTaskStatus(stateDir, name, taskCfg)
	if err := taskflow.UpdateTaskState(stateDir, name, taskCfg, status, end); err != nil {
		return 2, err
	}
	if runErr := taskflow.WriteTaskRun(stateDir, name, "git.branch.create", start, end, status, exitCodeFromError(err), strings.TrimSpace(string(output))); runErr != nil {
		logger.Error("failed to write git branch run", "task", name, "error", runErr)
	}

	if err != nil {
		return 2, fmt.Errorf("git branch create failed: %w", err)
	}
	return 0, nil
}

func resolveTaskStatus(stateDir string, name string, taskCfg config.TaskConfig) string {
	path := filepath.Join(stateDir, "tasks", name+".json")
	if record, err := state.ReadTask(path); err == nil && record.Status != "" {
		return record.Status
	}
	if taskCfg.Status != "" {
		return taskCfg.Status
	}
	return "planned"
}

func buildBranchName(taskName string, repo string) string {
	suffix := sanitizeTaskName(taskName)
	if repo == "external" {
		return fmt.Sprintf("task/external-%s", suffix)
	}
	return fmt.Sprintf("task/%s", suffix)
}

var branchCleaner = regexp.MustCompile(`[^a-z0-9._-]+`)

func sanitizeTaskName(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	clean := branchCleaner.ReplaceAllString(lower, "-")
	clean = strings.Trim(clean, "-")
	if clean == "" {
		return "task"
	}
	return clean
}
