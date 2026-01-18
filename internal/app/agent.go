package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"orchastration/internal/agent"
	"orchastration/internal/config"
	"orchastration/internal/logging"
)

func runAgent(args []string, _ config.Config, _ *logging.Logger, _ string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("agent requires a subcommand")
	}

	sub := args[0]
	switch sub {
	case "list":
		return agentList(os.Stdout)
	default:
		return 2, fmt.Errorf("unknown agent subcommand: %s", sub)
	}
}

func agentList(w io.Writer) (int, error) {
	return agentListWith(w, agent.List())
}

func agentListWith(w io.Writer, infos []agent.Info) (int, error) {
	if len(infos) == 0 {
		fmt.Fprintln(w, "no agents registered")
		return 0, nil
	}

	fmt.Fprintln(w, "available agents:")
	for _, info := range infos {
		caps := strings.Join(info.Capabilities, "; ")
		if caps == "" {
			caps = "(no capabilities listed)"
		}
		fmt.Fprintf(w, "- %s: %s\n", info.Name, caps)
	}
	return 0, nil
}
