# Orchastration

Orchastration is a cross-platform CLI that orchestrates deterministic, auditable execution of user-defined jobs and workflow tasks on a single machine. It also provides a file hashing command for integrity checks.

## Contents
- [How Orchastration Works (Diagrams)](#how-orchastration-works)
- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [Configuration](#configuration)
- [State](#state)
- [Logging](#logging)
- [Permissions](#permissions)

## How Orchastration Works

![Orchastration Internal Architecture v2](docs/diagrams/Orchastration%20F-Internal%20Architecture%20v2.svg)

Orchastration acts as a local workflow engine that coordinates task planning, execution, review, and documentation without taking ownership of target repositories. The orchestration engine drives the Planner, Builder, Reviewer, and Doc agents while preserving state and logs in OS-native directories. This separation keeps orchestration logic centralized while allowing domain logic to remain in the external project.

![Multi-Agent Workflow](docs/diagrams/Orchastration%20E-Multi-Agent%20Workflow.svg)

This workflow represents the multi-agent pipeline that powers v2: Planner decomposes goals, Builder executes tasks, Reviewer validates outputs, and Doc captures results. A shared OrchContext carries the plan, outputs, review results, and documentation references across stages.

![Context Sharing in Parallel](docs/diagrams/Orchastration%20G-Context%20Sharing%20%28Parallel%29.svg)

Parallel agent groups can run concurrently when configured, using the shared context to fan-out and fan-in results safely.

![External Repository Interaction](docs/diagrams/Orchastration%20H-External%20Repository%20Interaction.svg)

Orchastration executes commands in external repositories while keeping orchestration state and logs in its own directories. Outputs and documentation stay in the target repo; Orchastration only records metadata and run state.

## Features
- Native binaries for Windows and Linux
- Structured logging to stdout and a log file
- OS-appropriate config and log locations
- Deterministic, auditable job and task execution with local state records
- Multi-agent orchestration with sequential or parallel agent pipelines

## Install

### Linux
1. Build the binary:
   ```bash
   ./scripts/build-linux.sh
   ```
2. Run:
   ```bash
   ./dist/orchastration --help
   ```

Optional packaging:
```bash
./scripts/package-linux.sh
```
This produces `dist/orchastration-<version>-linux-amd64.tar.gz` and, if `dpkg-deb` is available, `dist/orchastration_<version>_amd64.deb`.

### Windows
1. Build the binary:
   ```powershell
   .\scripts\build-windows.ps1
   ```
2. Run:
   ```powershell
   .\dist\orchastration.exe --help
   ```

Optional packaging:
```powershell
.\scripts\package-windows.ps1
```
This produces `dist\orchastration-<version>-windows-amd64.zip`.

## Usage
```bash
orchastration --help
orchastration --version
orchastration hash --file ./path/to/file
orchastration list
orchastration run sample
orchastration status
orchastration plan list
orchastration plan create sample_task
orchastration plan status sample_task
orchastration build run sample_task
orchastration doc generate sample_task
orchastration git issue create sample_task
orchastration git branch create sample_task
orchastration agent list
orchastration orchestration list
orchastration orchestration run hello_multi_agent --goal "Implement login endpoint"
```

For a step-by-step walkthrough, see `USAGE.md`.

## Configuration
The config file is optional. If missing, defaults are used.
- Linux: `$XDG_CONFIG_HOME/orchastration/config.toml` (falls back to `~/.config/orchastration/config.toml`)
- Windows: `%AppData%\orchastration\config.toml`

Example config: `configs/config.example.toml`

## State
Execution records are stored under an OS-appropriate state directory (override with `--state-dir`):
- Linux: `$XDG_CACHE_HOME/orchastration/state` (falls back to `~/.cache/orchastration/state`)
- Windows: `%LocalAppData%\\orchastration\\state`

Job runs are recorded under:
`state/runs/<job-name>/<timestamp>.json` and `state/runs/<job-name>/last.json`

Task state is recorded under:
`state/tasks/<task>.json` and task runs under `state/runs/<task>/<timestamp>.json`

Orchestration runs are recorded under:
`state/orchestrations/<name>/<timestamp>.json`

## Logging
Logs are JSON and written to stdout and a log file:
- Linux: `$XDG_CACHE_HOME/orchastration/orchastration.log` (falls back to `~/.cache/orchastration/orchastration.log`)
- Windows: `%LocalAppData%\\orchastration\\orchastration.log`

## Permissions
Runs as a normal user. It only needs read access to files you hash and write access to your user config/cache directories.

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See LICENSE and NOTICE.
