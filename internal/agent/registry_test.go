package agent

import "testing"

type stubAgent struct {
	name string
	caps []string
}

func (a *stubAgent) Name() string {
	return a.name
}

func (a *stubAgent) Capabilities() []string {
	return a.caps
}

func (a *stubAgent) Execute(ctx *OrchContext) error {
	return nil
}

func TestRegistryRegisterListAndNew(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register("StubAgent", func() Agent {
		return &stubAgent{name: "StubAgent", caps: []string{"test capability"}}
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	list := reg.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(list))
	}
	if list[0].Name != "StubAgent" {
		t.Fatalf("expected name StubAgent, got %s", list[0].Name)
	}
	if len(list[0].Capabilities) != 1 || list[0].Capabilities[0] != "test capability" {
		t.Fatalf("unexpected capabilities: %#v", list[0].Capabilities)
	}

	instance, ok := reg.New("StubAgent")
	if !ok || instance == nil {
		t.Fatalf("expected New to return agent")
	}
	if instance.Name() != "StubAgent" {
		t.Fatalf("expected agent name StubAgent, got %s", instance.Name())
	}
}

func TestRegistryRejectsDuplicate(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register("DupAgent", func() Agent { return &stubAgent{name: "DupAgent"} }); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register("DupAgent", func() Agent { return &stubAgent{name: "DupAgent"} }); err == nil {
		t.Fatalf("expected duplicate registration error")
	}
}

func TestRegistryListOrder(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register("B-Agent", func() Agent { return &stubAgent{name: "B-Agent"} }); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register("A-Agent", func() Agent { return &stubAgent{name: "A-Agent"} }); err != nil {
		t.Fatalf("register: %v", err)
	}

	list := reg.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(list))
	}
	if list[0].Name != "A-Agent" || list[1].Name != "B-Agent" {
		t.Fatalf("unexpected order: %#v", list)
	}
}
