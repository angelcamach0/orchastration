package agent

// DocAgent documents results and outcomes.
type DocAgent struct{}

func (a *DocAgent) Name() string {
	return "DocAgent"
}

func (a *DocAgent) Capabilities() []string {
	return []string{"Document results and outcomes"}
}

func (a *DocAgent) Execute(ctx *OrchContext) error {
	return nil
}
