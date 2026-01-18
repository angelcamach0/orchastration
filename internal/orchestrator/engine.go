package orchestrator

import (
	"errors"
	"fmt"
	"path/filepath"
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

// Run executes agents sequentially for the named orchestration.
func (e *OrchestrationEngine) Run(orchestration string, agentNames []string, stateDir string) error {
	if orchestration == "" {
		return errors.New("orchestration name is required")
	}
	if len(agentNames) == 0 {
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
	agentRuns := make([]state.AgentRunRecord, 0, len(agentNames))
	status := "success"
	var execErr error

	for _, name := range agentNames {
		// TODO(v2): support parallel agent groups for concurrent orchestration.
		instance, ok := registry.New(name)
		if !ok {
			execErr = fmt.Errorf("agent not registered: %s", name)
			failedAt := e.now().UTC()
			agentRuns = append(agentRuns, state.AgentRunRecord{
				Name:       name,
				StartTime:  failedAt.Format(time.RFC3339),
				EndTime:    failedAt.Format(time.RFC3339),
				DurationMs: 0,
				Status:     "failed",
				Message:    execErr.Error(),
			})
			status = "failed"
			break
		}

		agentStart := e.now().UTC()
		runErr := instance.Execute(ctx)
		agentEnd := e.now().UTC()
		agentStatus := "success"
		message := "completed"
		if runErr != nil {
			agentStatus = "failed"
			message = runErr.Error()
			status = "failed"
			execErr = runErr
		}

		agentRuns = append(agentRuns, state.AgentRunRecord{
			Name:       name,
			StartTime:  agentStart.Format(time.RFC3339),
			EndTime:    agentEnd.Format(time.RFC3339),
			DurationMs: agentEnd.Sub(agentStart).Milliseconds(),
			Status:     agentStatus,
			Message:    message,
		})

		if runErr != nil {
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
