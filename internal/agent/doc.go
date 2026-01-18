package agent

import (
	"path/filepath"

	"orchastration/internal/taskflow"
)

// DocAgent documents results and outcomes.
type DocAgent struct{}

func (a *DocAgent) Name() string {
	return "DocAgent"
}

func (a *DocAgent) Capabilities() []string {
	return []string{"Document results and outcomes"}
}

func (a *DocAgent) Execute(ctx *OrchContext) error {
	tasks := readTaskList(ctx)
	if len(tasks) == 0 {
		return nil
	}

	deps, err := depsFromContext(ctx)
	if err != nil {
		return err
	}

	paths := make([]string, 0, len(tasks))
	for _, name := range tasks {
		if _, err := taskflow.DocGenerate(name, deps.cfg, deps.logger, deps.stateDir, deps.writer); err != nil {
			return err
		}
		taskCfg, ok := deps.cfg.Tasks[name]
		if !ok {
			continue
		}
		docPath := filepath.Join(taskCfg.WorkingDir, "docs", "tasks", name+".md")
		paths = append(paths, docPath)
	}
	if len(paths) > 0 {
		ctx.Set(ctxKeyDocPaths, paths)
	}
	return nil
}
