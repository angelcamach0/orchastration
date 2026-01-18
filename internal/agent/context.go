package agent

import "sync"

// OrchContext holds shared orchestration state for a run.
type OrchContext struct {
	mu   sync.RWMutex
	data map[string]any
}

// Set stores a value in the context.
func (c *OrchContext) Set(key string, value any) {
	if c == nil {
		return
	}
	c.mu.Lock()
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data[key] = value
	c.mu.Unlock()
}

// Get returns a value from the context.
func (c *OrchContext) Get(key string) (any, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.data == nil {
		return nil, false
	}
	val, ok := c.data[key]
	return val, ok
}

// SnapshotStrings returns a copy of string values for persistence.
func (c *OrchContext) SnapshotStrings() map[string]string {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.data == nil {
		return nil
	}
	out := make(map[string]string, len(c.data))
	for key, value := range c.data {
		str, ok := value.(string)
		if !ok {
			continue
		}
		out[key] = str
	}
	return out
}
