package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/state"
)

func runBuild(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("build requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "run":
		return buildRun(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown build subcommand: %s", sub)
	}
}

func buildRun(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("build run requires a task name")
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
	if err := updateTaskState(stateDir, name, taskCfg, "in_progress", start); err != nil {
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

	if err := updateTaskState(stateDir, name, taskCfg, status, end); err != nil {
		return 2, err
	}
	if err := writeTaskRun(stateDir, name, "build.run", start, end, status, exitCode, message); err != nil {
		logger.Error("failed to write build run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(os.Stdout, "task=%s exit=%d status=%s\n", name, exitCode, status)
	if execErr != nil {
		return 2, execErr
	}
	return 0, nil
}

func updateTaskState(stateDir string, name string, taskCfg config.TaskConfig, status string, lastRun time.Time) error {
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
