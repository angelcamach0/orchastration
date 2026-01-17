package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/state"
)

func runPlan(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("plan requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "list":
		return planList(cfg, stateDir)
	case "create":
		return planCreate(args[1:], cfg, logger, stateDir)
	case "status":
		return planStatus(args[1:], cfg, stateDir)
	default:
		return 2, fmt.Errorf("unknown plan subcommand: %s", sub)
	}
}

func planList(cfg config.Config, stateDir string) (int, error) {
	if len(cfg.Tasks) == 0 {
		fmt.Fprintln(os.Stdout, "no tasks configured")
		return 0, nil
	}

	names := make([]string, 0, len(cfg.Tasks))
	for name := range cfg.Tasks {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		status := cfg.Tasks[name].Status
		if status == "" {
			status = "planned"
		}
		taskPath := filepath.Join(stateDir, "tasks", name+".json")
		record, err := state.ReadTask(taskPath)
		if err == nil && record.Status != "" {
			status = record.Status
		}
		line := name
		if desc := cfg.Tasks[name].Description; desc != "" {
			line = fmt.Sprintf("%s - %s", name, desc)
		}
		fmt.Fprintf(os.Stdout, "%s [%s]\n", line, status)
	}

	return 0, nil
}

func planCreate(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("plan create requires a task name")
	}

	name := args[0]
	taskCfg, ok := cfg.Tasks[name]
	if !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}
	if err := validateTaskConfig(name, taskCfg); err != nil {
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

	if err := writeTaskRun(stateDir, name, "plan.create", now, now, status, 0, "created"); err != nil {
		logger.Error("failed to write plan run", "task", name, "error", err)
		return 2, err
	}

	fmt.Fprintf(os.Stdout, "task=%s status=%s\n", name, status)
	return 0, nil
}

func planStatus(args []string, cfg config.Config, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("plan status requires a task name")
	}

	name := args[0]
	if _, ok := cfg.Tasks[name]; !ok {
		return 2, fmt.Errorf("unknown task: %s", name)
	}

	taskPath := filepath.Join(stateDir, "tasks", name+".json")
	record, err := state.ReadTask(taskPath)
	if err != nil {
		return 2, fmt.Errorf("task not initialized: %s", name)
	}

	fmt.Fprintf(os.Stdout, "%s status=%s last_run=%s\n", name, record.Status, record.LastRun)
	return 0, nil
}

func writeTaskRun(stateDir string, taskName string, action string, start time.Time, end time.Time, status string, exitCode int, message string) error {
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

func validateTaskConfig(name string, task config.TaskConfig) error {
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
	if len(task.Command) == 0 {
		return fmt.Errorf("task %s has empty command", name)
	}
	return nil
}
