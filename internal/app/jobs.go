package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/state"
)

func listJobs(cfg config.Config) (int, error) {
	if len(cfg.Jobs) == 0 {
		fmt.Fprintln(os.Stdout, "no jobs configured")
		return 0, nil
	}

	names := make([]string, 0, len(cfg.Jobs))
	for name := range cfg.Jobs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		job := cfg.Jobs[name]
		line := name
		if job.Description != "" {
			line = fmt.Sprintf("%s - %s", name, job.Description)
		}
		fmt.Fprintln(os.Stdout, line)
	}
	return 0, nil
}

func runJob(args []string, cfg config.Config, logger *logging.Logger, stateDir string, version string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("run requires a job name")
	}

	jobName := args[0]
	job, ok := cfg.Jobs[jobName]
	if !ok {
		return 2, fmt.Errorf("unknown job: %s", jobName)
	}

	if len(job.Command) == 0 {
		return 2, fmt.Errorf("job %s has empty command", jobName)
	}

	start := time.Now().UTC()
	timeStamp := start.Format("20060102T150405Z")
	runDir := filepath.Join(stateDir, "runs", jobName)
	stdoutPath := filepath.Join(runDir, timeStamp+".stdout.log")
	stderrPath := filepath.Join(runDir, timeStamp+".stderr.log")

	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return 2, fmt.Errorf("create run dir: %w", err)
	}

	stdoutFile, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 2, fmt.Errorf("open stdout file: %w", err)
	}
	defer stdoutFile.Close()

	stderrFile, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 2, fmt.Errorf("open stderr file: %w", err)
	}
	defer stderrFile.Close()

	ctx := context.Background()
	if job.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, job.Command[0], job.Command[1:]...)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	if job.WorkingDir != "" {
		cmd.Dir = job.WorkingDir
	}
	cmd.Env = mergeEnv(job.Env)

	logger.Info("job starting", "job", jobName, "command", strings.Join(job.Command, " "))
	execErr := cmd.Run()

	end := time.Now().UTC()
	duration := end.Sub(start)
	exitCode := exitCodeFromError(execErr)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logger.Error("job timed out", "job", jobName, "timeout_seconds", job.TimeoutSeconds)
	}
	if execErr != nil {
		logger.Error("job failed", "job", jobName, "error", execErr)
	}

	record := state.Record{
		JobName:    jobName,
		StartTime:  start.Format(time.RFC3339),
		EndTime:    end.Format(time.RFC3339),
		DurationMs: duration.Milliseconds(),
		ExitCode:   exitCode,
		StdoutPath: stdoutPath,
		StderrPath: stderrPath,
		OS:         runtime.GOOS,
		Version:    version,
	}

	recordPath := filepath.Join(runDir, timeStamp+".json")
	if err := state.WriteRecord(recordPath, record); err != nil {
		logger.Error("failed to write record", "job", jobName, "error", err)
		return 2, err
	}
	lastPath := filepath.Join(runDir, "last.json")
	if err := state.WriteRecord(lastPath, record); err != nil {
		logger.Error("failed to write last record", "job", jobName, "error", err)
		return 2, err
	}

	fmt.Fprintf(os.Stdout, "job=%s exit=%d duration_ms=%d\n", jobName, exitCode, duration.Milliseconds())
	if execErr != nil {
		return 2, execErr
	}
	return 0, nil
}

func jobStatus(cfg config.Config, stateDir string) (int, error) {
	if len(cfg.Jobs) == 0 {
		fmt.Fprintln(os.Stdout, "no jobs configured")
		return 0, nil
	}

	names := make([]string, 0, len(cfg.Jobs))
	for name := range cfg.Jobs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		lastPath := filepath.Join(stateDir, "runs", name, "last.json")
		record, err := state.ReadRecord(lastPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s - no runs recorded\n", name)
			continue
		}
		fmt.Fprintf(os.Stdout, "%s - exit=%d duration_ms=%d start=%s\n", name, record.ExitCode, record.DurationMs, record.StartTime)
	}

	return 0, nil
}

func mergeEnv(extra map[string]string) []string {
	env := os.Environ()
	if len(extra) == 0 {
		return env
	}

	seen := make(map[string]struct{}, len(env))
	for i, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		key := parts[0]
		if value, ok := extra[key]; ok {
			env[i] = key + "=" + value
		}
		seen[key] = struct{}{}
	}

	for key, value := range extra {
		if _, ok := seen[key]; ok {
			continue
		}
		env = append(env, key+"="+value)
	}

	return env
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
