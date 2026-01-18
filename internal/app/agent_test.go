package app

import (
	"bytes"
	"testing"

	"orchastration/internal/agent"
)

func TestAgentListWithAgents(t *testing.T) {
	var buf bytes.Buffer
	infos := []agent.Info{
		{Name: "BuilderAgent", Capabilities: []string{"Build outputs"}},
		{Name: "PlannerAgent", Capabilities: []string{"Create plan"}},
	}

	if _, err := agentListWith(&buf, infos); err != nil {
		t.Fatalf("agentListWith: %v", err)
	}

	expected := "available agents:\n- BuilderAgent: Build outputs\n- PlannerAgent: Create plan\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}

func TestAgentListWithEmpty(t *testing.T) {
	var buf bytes.Buffer
	if _, err := agentListWith(&buf, nil); err != nil {
		t.Fatalf("agentListWith: %v", err)
	}

	expected := "no agents registered\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}
