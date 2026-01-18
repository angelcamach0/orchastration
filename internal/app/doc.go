package app

import (
	"errors"
	"fmt"
	"os"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/taskflow"
)

func runDoc(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("doc requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "generate":
		return docGenerate(args[1:], cfg, logger, stateDir)
	default:
		return 2, fmt.Errorf("unknown doc subcommand: %s", sub)
	}
}

func docGenerate(args []string, cfg config.Config, logger *logging.Logger, stateDir string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("doc generate requires a task name")
	}

	name := args[0]
	return taskflow.DocGenerate(name, cfg, logger, stateDir, os.Stdout)
}
