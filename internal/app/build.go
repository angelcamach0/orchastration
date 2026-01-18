package app

import (
	"errors"
	"fmt"
	"os"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/taskflow"
)

func runBuild(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("build requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "run":
		return buildRun(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown build subcommand: %s", sub)
	}
}

func buildRun(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("build run requires a task name")
	}

	name := args[0]
	return taskflow.BuildRun(name, cfg, logger, stateDir, os.Stdout)
}
