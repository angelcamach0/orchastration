package agent

import (
	"path/filepath"

	"orchastration/internal/taskflow"
)

// BuilderAgent executes a task plan to produce outputs.
type BuilderAgent struct{}

func (a *BuilderAgent) Name() string {
	return "BuilderAgent"
}

func (a *BuilderAgent) Capabilities() []string {
	return []string{"Execute planned tasks to produce outputs"}
}

func (a *BuilderAgent) Execute(ctx *OrchContext) error {
	tasks := readTaskList(ctx)
	if len(tasks) == 0 {
		return nil
	}

	deps, err := depsFromContext(ctx)
	if err != nil {
		return err
	}

	outputs := make([]string, 0)
	for _, name := range tasks {
		if _, err := taskflow.BuildRun(name, deps.cfg, deps.logger, deps.stateDir, deps.writer); err != nil {
			return err
		}

		taskCfg, ok := deps.cfg.Tasks[name]
		if !ok {
			continue
		}
		for _, output := range taskCfg.Outputs {
			if filepath.IsAbs(output) {
				outputs = append(outputs, output)
			} else {
				outputs = append(outputs, filepath.Join(taskCfg.WorkingDir, output))
			}
		}
	}
	if len(outputs) > 0 {
		ctx.Set(ctxKeyBuildOutputs, outputs)
	}
	return nil
}
