package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"orchastration/internal/config"
	"orchastration/internal/logging"
	"orchastration/internal/platform"
	"orchastration/internal/version"
)

const appName = "orchastration"

func Run(args []string, ver version.Info) (int, error) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return 0, nil
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printUsage(os.Stdout)
			return 0, nil
		}
		if arg == "--version" {
			fmt.Fprintln(os.Stdout, ver.String())
			return 0, nil
		}
	}

	root := flag.NewFlagSet(appName, flag.ContinueOnError)
	root.SetOutput(io.Discard)
	configPathFlag := root.String("config", "", "path to config file")
	stateDirFlag := root.String("state-dir", "", "path to state directory")
	if err := root.Parse(args); err != nil {
		return 2, err
	}

	remaining := root.Args()
	if len(remaining) == 0 {
		printUsage(os.Stdout)
		return 0, nil
	}

	cfgPath := *configPathFlag
	if cfgPath == "" {
		cfgPath = platform.DefaultConfigPath(appName)
	}
	stateDir := *stateDirFlag
	if stateDir == "" {
		stateDir = platform.DefaultStateDir(appName)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return 2, fmt.Errorf("load config: %w", err)
	}

	logPath := platform.DefaultLogPath(appName)
	logger, err := logging.New(cfg.Logging.Level, logPath)
	if err != nil {
		return 2, fmt.Errorf("init logger: %w", err)
	}

	cmd := remaining[0]
	switch cmd {
	case "hash":
		return runHash(remaining[1:], cfg, logger)
	case "run":
		return runJob(remaining[1:], cfg, logger, stateDir, ver.String())
	case "plan":
		return runPlan(remaining[1:], cfg, logger, stateDir)
	case "build":
		return runBuild(remaining[1:], cfg, logger, stateDir)
	case "doc":
		return runDoc(remaining[1:], cfg, logger, stateDir)
	case "git":
		return runGit(remaining[1:], cfg, logger, stateDir)
	case "agent":
		return runAgent(remaining[1:], cfg, logger, stateDir)
	case "orchestration":
		return runOrchestration(remaining[1:], cfg, logger, stateDir)
	case "list":
		return listJobs(cfg)
	case "status":
		return jobStatus(cfg, stateDir)
	default:
		return 2, fmt.Errorf("unknown command: %s", cmd)
	}
}

func runHash(args []string, cfg config.Config, logger *logging.Logger) (int, error) {
	fs := flag.NewFlagSet("hash", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	filePath := fs.String("file", "", "path to file")
	algo := fs.String("algo", cfg.Hash.Algorithm, "hash algorithm (sha256, sha1, sha512)")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}

	if *filePath == "" {
		return 2, errors.New("hash requires --file")
	}

	absPath, err := filepath.Abs(*filePath)
	if err != nil {
		return 2, fmt.Errorf("resolve path: %w", err)
	}

	result, err := hashFile(absPath, *algo)
	if err != nil {
		logger.Error("hash failed", "file", absPath, "algorithm", *algo, "error", err)
		return 2, err
	}

	logger.Info("hash computed", "file", absPath, "algorithm", *algo)
	fmt.Fprintf(os.Stdout, "{\"file\":%q,\"algorithm\":%q,\"hash\":%q}\n", absPath, *algo, result)
	return 0, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "orchastration - cross-platform orchestration helper")
	fmt.Fprintln(w, "\nUsage:")
	fmt.Fprintln(w, "  orchastration [--config path] [--state-dir path] <command> [options]")
	fmt.Fprintln(w, "\nCommands:")
	fmt.Fprintln(w, "  hash   Compute file hash (useful for integrity checks)")
	fmt.Fprintln(w, "  run    Run a configured job by name")
	fmt.Fprintln(w, "  list   List configured jobs")
	fmt.Fprintln(w, "  status Show last recorded job runs")
	fmt.Fprintln(w, "  plan   Plan workflow tasks (list, create, status)")
	fmt.Fprintln(w, "  build  Run workflow tasks")
	fmt.Fprintln(w, "  doc    Generate task documentation")
	fmt.Fprintln(w, "  git    Git helpers (issue, branch)")
	fmt.Fprintln(w, "  agent  Agent helpers (list)")
	fmt.Fprintln(w, "  orchestration  Orchestration helpers (list, run)")
	fmt.Fprintln(w, "\nFlags:")
	fmt.Fprintln(w, "  --help       Show help")
	fmt.Fprintln(w, "  --version    Show version")
	fmt.Fprintln(w, "  --config     Path to config file")
	fmt.Fprintln(w, "  --state-dir  Path to state directory")
}
