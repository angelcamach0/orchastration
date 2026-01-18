package orchestrator

import (
	"os"
	"path/filepath"
	"testing"

	"orchastration/internal/agent"
	"orchastration/internal/state"
)

type recordingAgent struct {
	name  string
	order *[]string
}

func (a *recordingAgent) Name() string {
	return a.name
}

func (a *recordingAgent) Capabilities() []string {
	return []string{"record execution order"}
}

func (a *recordingAgent) Execute(ctx *agent.OrchContext) error {
	*a.order = append(*a.order, a.name)
	return nil
}

func TestEngineRunSequentialOrder(t *testing.T) {
	reg := agent.NewRegistry()
	order := []string{}
	if err := reg.Register("PlannerAgent", func() agent.Agent {
		return &recordingAgent{name: "PlannerAgent", order: &order}
	}); err != nil {
		t.Fatalf("register planner: %v", err)
	}
	if err := reg.Register("BuilderAgent", func() agent.Agent {
		return &recordingAgent{name: "BuilderAgent", order: &order}
	}); err != nil {
		t.Fatalf("register builder: %v", err)
	}
	if err := reg.Register("ReviewerAgent", func() agent.Agent {
		return &recordingAgent{name: "ReviewerAgent", order: &order}
	}); err != nil {
		t.Fatalf("register reviewer: %v", err)
	}
	if err := reg.Register("DocAgent", func() agent.Agent {
		return &recordingAgent{name: "DocAgent", order: &order}
	}); err != nil {
		t.Fatalf("register doc: %v", err)
	}

	engine := NewOrchestrationEngine(reg)
	stateDir := t.TempDir()

	if err := engine.Run("pipeline", []string{
		"PlannerAgent",
		"BuilderAgent",
		"ReviewerAgent",
		"DocAgent",
	}, stateDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	expectedOrder := []string{"PlannerAgent", "BuilderAgent", "ReviewerAgent", "DocAgent"}
	if len(order) != len(expectedOrder) {
		t.Fatalf("expected %d agent calls, got %d", len(expectedOrder), len(order))
	}
	for i, name := range expectedOrder {
		if order[i] != name {
			t.Fatalf("expected order %v, got %v", expectedOrder, order)
		}
	}

	runDir := filepath.Join(stateDir, "orchestrations", "pipeline")
	entries, err := os.ReadDir(runDir)
	if err != nil {
		t.Fatalf("read run dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 run record, got %d", len(entries))
	}

	recordPath := filepath.Join(runDir, entries[0].Name())
	record, err := state.ReadOrchestrationRun(recordPath)
	if err != nil {
		t.Fatalf("read record: %v", err)
	}
	if record.Orchestration != "pipeline" {
		t.Fatalf("unexpected orchestration: %s", record.Orchestration)
	}
	if record.Status != "success" {
		t.Fatalf("unexpected status: %s", record.Status)
	}
	if len(record.Agents) != 4 {
		t.Fatalf("expected 4 agent records, got %d", len(record.Agents))
	}
	for i, name := range expectedOrder {
		if record.Agents[i].Name != name {
			t.Fatalf("unexpected agent record order: %v", record.Agents)
		}
	}
}
