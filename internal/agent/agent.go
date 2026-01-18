package agent

// Agent defines the common contract for all orchestration agents.
type Agent interface {
	Name() string
	Capabilities() []string
	Execute(ctx *OrchContext) error
}
