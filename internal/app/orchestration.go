package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/orchestrator"
)

func runOrchestration(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("orchestration requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "list":
		return orchestrationList(os.Stdout, cfg.Orchestrations)
	case "run":
		return orchestrationRun(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown orchestration subcommand: %s", sub)
	}
}

func orchestrationList(w io.Writer, orchestrations map[string]config.OrchestrationConfig) (int, error) {
	if len(orchestrations) == 0 {
		fmt.Fprintln(w, "no orchestrations configured")
		return 0, nil
	}

	names := make([]string, 0, len(orchestrations))
	for name := range orchestrations {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintln(w, "available orchestrations:")
	for _, name := range names {
		desc := orchestrations[name].Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Fprintf(w, "- %s: %s\n", name, desc)
	}
	return 0, nil
}

func orchestrationRun(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) < 1 {
		return 2, errors.New("orchestration run requires a name")
	}

	name := args[0]
	orchCfg, ok := cfg.Orchestrations[name]
	if !ok {
		return 2, fmt.Errorf("unknown orchestration: %s", name)
	}
	steps := buildOrchestrationSteps(orchCfg)
	if len(steps) == 0 {
		return 2, fmt.Errorf("orchestration %s has no agents", name)
	}

	engine := orchestrator.NewOrchestrationEngine(nil)
	if err := engine.Run(name, steps, stateDir); err != nil {
		logger.Error("orchestration failed", "orchestration", name, "error", err)
		return 2, err
	}
	return 0, nil
}

func buildOrchestrationSteps(cfg config.OrchestrationConfig) [][]string {
	if len(cfg.Steps) != 0 {
		return cfg.Steps
	}
	steps := make([][]string, 0, len(cfg.Agents))
	for _, name := range cfg.Agents {
		steps = append(steps, []string{name})
	}
	return steps
}
