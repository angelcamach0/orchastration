package orchestrator

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"orchastration/internal/agent"
	"orchastration/internal/state"
)

// OrchestrationEngine coordinates agent execution for a run.
type OrchestrationEngine struct {
	Registry *agent.Registry
	now      func() time.Time
}

// NewOrchestrationEngine builds an engine with the provided registry.
// If registry is nil, the default agent registry is used.
func NewOrchestrationEngine(registry *agent.Registry) *OrchestrationEngine {
	if registry == nil {
		registry = agent.DefaultRegistry()
	}
	return &OrchestrationEngine{
		Registry: registry,
		now:      time.Now,
	}
}

// Run executes agents sequentially or concurrently based on step groups.
func (e *OrchestrationEngine) Run(orchestration string, steps [][]string, stateDir string) error {
	if orchestration == "" {
		return errors.New("orchestration name is required")
	}
	if len(steps) == 0 {
		return errors.New("at least one agent is required")
	}
	if stateDir == "" {
		return errors.New("state directory is required")
	}
	if e == nil {
		return errors.New("orchestration engine is nil")
	}

	registry := e.Registry
	if registry == nil {
		registry = agent.DefaultRegistry()
	}

	ctx := &agent.OrchContext{}
	start := e.now().UTC()
	agentRuns := make([]state.AgentRunRecord, 0, len(steps))
	status := "success"
	var execErr error

	for _, group := range steps {
		if len(group) == 0 {
			continue
		}

		if len(group) == 1 {
			runRecord, runErr := runAgent(ctx, registry, e.now, group[0])
			agentRuns = append(agentRuns, runRecord)
			if runErr != nil {
				status = "failed"
				execErr = runErr
				break
			}
			continue
		}

		groupRuns, groupErr := runParallelGroup(ctx, registry, e.now, group)
		agentRuns = append(agentRuns, groupRuns...)
		if groupErr != nil {
			status = "failed"
			execErr = groupErr
			break
		}
	}

	end := e.now().UTC()
	record := state.OrchestrationRunRecord{
		Orchestration: orchestration,
		Agents:        agentRuns,
		StartTime:     start.Format(time.RFC3339),
		EndTime:       end.Format(time.RFC3339),
		DurationMs:    end.Sub(start).Milliseconds(),
		Status:        status,
		Context:       ctx.SnapshotStrings(),
	}

	writeErr := writeOrchestrationRun(stateDir, orchestration, start, record)
	if execErr != nil && writeErr != nil {
		return errors.Join(execErr, writeErr)
	}
	if writeErr != nil {
		return writeErr
	}
	return execErr
}

func writeOrchestrationRun(stateDir string, orchestration string, start time.Time, record state.OrchestrationRunRecord) error {
	timestamp := start.Format("20060102T150405Z")
	runPath := filepath.Join(stateDir, "orchestrations", orchestration, timestamp+".json")
	return state.WriteOrchestrationRun(runPath, record)
}

func runAgent(ctx *agent.OrchContext, registry *agent.Registry, now func() time.Time, name string) (state.AgentRunRecord, error) {
	instance, ok := registry.New(name)
	if !ok {
		failedAt := now().UTC()
		err := fmt.Errorf("agent not registered: %s", name)
		return state.AgentRunRecord{
			Name:       name,
			StartTime:  failedAt.Format(time.RFC3339),
			EndTime:    failedAt.Format(time.RFC3339),
			DurationMs: 0,
			Status:     "failed",
			Message:    err.Error(),
		}, err
	}

	start := now().UTC()
	runErr := instance.Execute(ctx)
	end := now().UTC()
	status := "success"
	message := "completed"
	if runErr != nil {
		status = "failed"
		message = runErr.Error()
	}
	return state.AgentRunRecord{
		Name:       name,
		StartTime:  start.Format(time.RFC3339),
		EndTime:    end.Format(time.RFC3339),
		DurationMs: end.Sub(start).Milliseconds(),
		Status:     status,
		Message:    message,
	}, runErr
}

func runParallelGroup(ctx *agent.OrchContext, registry *agent.Registry, now func() time.Time, names []string) ([]state.AgentRunRecord, error) {
	results := make([]state.AgentRunRecord, len(names))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, name := range names {
		instance, ok := registry.New(name)
		if !ok {
			failedAt := now().UTC()
			err := fmt.Errorf("agent not registered: %s", name)
			results[i] = state.AgentRunRecord{
				Name:       name,
				StartTime:  failedAt.Format(time.RFC3339),
				EndTime:    failedAt.Format(time.RFC3339),
				DurationMs: 0,
				Status:     "failed",
				Message:    err.Error(),
			}
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		wg.Add(1)
		go func(idx int, agentName string, agentInstance agent.Agent) {
			defer wg.Done()
			start := now().UTC()
			runErr := agentInstance.Execute(ctx)
			end := now().UTC()
			status := "success"
			message := "completed"
			if runErr != nil {
				status = "failed"
				message = runErr.Error()
			}

			results[idx] = state.AgentRunRecord{
				Name:       agentName,
				StartTime:  start.Format(time.RFC3339),
				EndTime:    end.Format(time.RFC3339),
				DurationMs: end.Sub(start).Milliseconds(),
				Status:     status,
				Message:    message,
			}

			if runErr != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = runErr
				}
				mu.Unlock()
			}
		}(i, name, instance)
	}

	wg.Wait()
	return results, firstErr
}
