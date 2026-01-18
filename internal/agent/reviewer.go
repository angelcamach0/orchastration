package agent

// ReviewerAgent validates outputs for quality and correctness.
type ReviewerAgent struct{}

func (a *ReviewerAgent) Name() string {
	return "ReviewerAgent"
}

func (a *ReviewerAgent) Capabilities() []string {
	return []string{"Review outputs for quality and correctness"}
}

func (a *ReviewerAgent) Execute(ctx *OrchContext) error {
	return nil
}
