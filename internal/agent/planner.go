package agent

// PlannerAgent decomposes a high-level goal into a task plan.
type PlannerAgent struct{}

func (a *PlannerAgent) Name() string {
	return "PlannerAgent"
}

func (a *PlannerAgent) Capabilities() []string {
	return []string{"Create a structured task plan from a goal"}
}

func (a *PlannerAgent) Execute(ctx *OrchContext) error {
	return nil
}
