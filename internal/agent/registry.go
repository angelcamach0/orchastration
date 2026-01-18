package agent

import (
	"fmt"
	"sort"
	"sync"
)

// Constructor builds a new Agent instance.
type Constructor func() Agent

// Info describes a registered agent.
type Info struct {
	Name         string
	Capabilities []string
}

// Registry stores agent constructors by name.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Constructor
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]Constructor)}
}

// Register adds a constructor to the registry.
func (r *Registry) Register(name string, ctor Constructor) error {
	if name == "" {
		return fmt.Errorf("agent name is required")
	}
	if ctor == nil {
		return fmt.Errorf("agent constructor is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("agent %s already registered", name)
	}
	r.entries[name] = ctor
	return nil
}

// New returns a new agent instance by name.
func (r *Registry) New(name string) (Agent, bool) {
	r.mu.RLock()
	ctor, ok := r.entries[name]
	r.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return ctor(), true
}

// List returns the registered agent info in stable order.
func (r *Registry) List() []Info {
	r.mu.RLock()
	names := make([]string, 0, len(r.entries))
	for name := range r.entries {
		names = append(names, name)
	}
	r.mu.RUnlock()

	sort.Strings(names)
	infos := make([]Info, 0, len(names))
	for _, name := range names {
		r.mu.RLock()
		ctor := r.entries[name]
		r.mu.RUnlock()
		caps := []string{}
		if ctor != nil {
			if agent := ctor(); agent != nil {
				caps = agent.Capabilities()
			}
		}
		infos = append(infos, Info{Name: name, Capabilities: caps})
	}
	return infos
}

var defaultRegistry = NewRegistry()

// Register adds a constructor to the default registry.
func Register(name string, ctor Constructor) error {
	return defaultRegistry.Register(name, ctor)
}

// MustRegister adds a constructor or panics if registration fails.
func MustRegister(name string, ctor Constructor) {
	if err := Register(name, ctor); err != nil {
		panic(err)
	}
}

// New returns a new agent instance from the default registry.
func New(name string) (Agent, bool) {
	return defaultRegistry.New(name)
}

// List returns registered agent info from the default registry.
func List() []Info {
	return defaultRegistry.List()
}

func init() {
	MustRegister("PlannerAgent", func() Agent { return &PlannerAgent{} })
	MustRegister("BuilderAgent", func() Agent { return &BuilderAgent{} })
	MustRegister("ReviewerAgent", func() Agent { return &ReviewerAgent{} })
	MustRegister("DocAgent", func() Agent { return &DocAgent{} })
}
