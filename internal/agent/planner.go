package agent

import "orchastration/internal/taskflow"

// PlannerAgent decomposes a high-level goal into a task plan.
type PlannerAgent struct{}

func (a *PlannerAgent) Name() string {
	return "PlannerAgent"
}

func (a *PlannerAgent) Capabilities() []string {
	return []string{"Create a structured task plan from a goal"}
}

func (a *PlannerAgent) Execute(ctx *OrchContext) error {
	if ctx == nil {
		return nil
	}

	if goalValue, ok := ctx.Get(ctxKeyGoal); ok {
		if goal, ok := goalValue.(string); ok && goal != "" {
			ctx.Set(ctxKeyPlanGoal, goal)
		}
	}

	taskName := readTaskName(ctx)
	if taskName == "" {
		ctx.Set(ctxKeyPlanTasks, []string{})
		return nil
	}

	ctx.Set(ctxKeyPlanTasks, []string{taskName})

	deps, err := depsFromContext(ctx)
	if err != nil {
		return err
	}
	_, err = taskflow.PlanCreate(taskName, deps.cfg, deps.logger, deps.stateDir, deps.writer)
	return err
}
