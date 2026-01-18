package agent

import (
	"errors"
	"io"

	"orchastration/internal/config"
	"orchastration/internal/logging"
)

const (
	ctxKeyConfig   = "config"
	ctxKeyLogger   = "logger"
	ctxKeyStateDir = "state.dir"
	ctxKeyWriter   = "writer"
	ctxKeyGoal     = "goal"
	ctxKeyTaskName = "task.name"
	ctxKeyPlanGoal = "plan.goal"
	ctxKeyPlanTasks = "plan.tasks"
	ctxKeyBuildOutputs = "build.outputs"
	ctxKeyReviewStatus = "review.status"
	ctxKeyReviewReport = "review.report"
	ctxKeyDocPaths = "doc.paths"
)

type deps struct {
	cfg      config.Config
	logger   *logging.Logger
	stateDir string
	writer   io.Writer
}

func depsFromContext(ctx *OrchContext) (deps, error) {
	if ctx == nil {
		return deps{}, errors.New("context is nil")
	}

	rawCfg, ok := ctx.Get(ctxKeyConfig)
	if !ok {
		return deps{}, errors.New("missing config in context")
	}
	cfg, ok := rawCfg.(config.Config)
	if !ok {
		return deps{}, errors.New("invalid config type in context")
	}

	rawLogger, ok := ctx.Get(ctxKeyLogger)
	if !ok {
		return deps{}, errors.New("missing logger in context")
	}
	logger, ok := rawLogger.(*logging.Logger)
	if !ok {
		return deps{}, errors.New("invalid logger type in context")
	}

	rawStateDir, ok := ctx.Get(ctxKeyStateDir)
	if !ok {
		return deps{}, errors.New("missing state dir in context")
	}
	stateDir, ok := rawStateDir.(string)
	if !ok {
		return deps{}, errors.New("invalid state dir type in context")
	}

	rawWriter, ok := ctx.Get(ctxKeyWriter)
	if !ok {
		return deps{}, errors.New("missing writer in context")
	}
	writer, ok := rawWriter.(io.Writer)
	if !ok {
		return deps{}, errors.New("invalid writer type in context")
	}

	return deps{
		cfg:      cfg,
		logger:   logger,
		stateDir: stateDir,
		writer:   writer,
	}, nil
}

func readTaskName(ctx *OrchContext) string {
	if ctx == nil {
		return ""
	}
	if value, ok := ctx.Get(ctxKeyTaskName); ok {
		if name, ok := value.(string); ok {
			return name
		}
	}
	return ""
}

func readTaskList(ctx *OrchContext) []string {
	if ctx == nil {
		return nil
	}
	if value, ok := ctx.Get(ctxKeyPlanTasks); ok {
		if names, ok := value.([]string); ok {
			return names
		}
	}
	if name := readTaskName(ctx); name != "" {
		return []string{name}
	}
	return nil
}
