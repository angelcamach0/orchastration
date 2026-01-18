package app

import (
	"bytes"
	"testing"

	"orchastration/internal/config"
)

func TestOrchestrationList(t *testing.T) {
	orchestrations := map[string]config.OrchestrationConfig{
		"beta":  {Description: "Second"},
		"alpha": {Description: "First"},
	}

	var buf bytes.Buffer
	if _, err := orchestrationList(&buf, orchestrations); err != nil {
		t.Fatalf("orchestrationList: %v", err)
	}

	expected := "available orchestrations:\n- alpha: First\n- beta: Second\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}

func TestOrchestrationListEmpty(t *testing.T) {
	var buf bytes.Buffer
	if _, err := orchestrationList(&buf, nil); err != nil {
		t.Fatalf("orchestrationList: %v", err)
	}

	expected := "no orchestrations configured\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}
