package agent

// BuilderAgent executes a task plan to produce outputs.
type BuilderAgent struct{}

func (a *BuilderAgent) Name() string {
	return "BuilderAgent"
}

func (a *BuilderAgent) Capabilities() []string {
	return []string{"Execute planned tasks to produce outputs"}
}

func (a *BuilderAgent) Execute(ctx *OrchContext) error {
	return nil
}
