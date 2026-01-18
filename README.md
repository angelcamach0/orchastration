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

![Orchastration Unified System Diagram](docs/diagrams/Orchastration%200-Unified%20System%20Diagram.svg)

Orchastration acts as a local workflow agent that coordinates task planning, execution, and documentation without taking ownership of the target repositories. The diagram highlights how Orchastration stays separate from external repos, running commands inside them while keeping state and logs in its own directories. At a high level, Planner, Builder, and Documentor work together to make task intent explicit, run deterministic commands, and capture outcomes. This separation keeps orchestration logic centralized while allowing domain logic to remain in the external project.

![Planner Builder Documentor Loop](docs/diagrams/Orchastration%20C-Planner%20%E2%86%92%20Builder%20%E2%86%92%20Documentor%20Loop%20%28Core%20Value%29.svg)

This loop represents the repeatable lifecycle of a task from plan to build to documentation. Planning establishes intent before any command runs, building executes the scoped command with explicit inputs, and documentation records results as a first-class output. The cycle is designed to be auditable and repeatable so that outcomes can be recreated and reviewed. Documentation is not an afterthought; it is the final, required step in the loop.

## Features
- Native binaries for Windows and Linux
- Structured logging to stdout and a log file
- OS-appropriate config and log locations
- Deterministic, auditable job and task execution with local state records

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

## Logging
Logs are JSON and written to stdout and a log file:
- Linux: `$XDG_CACHE_HOME/orchastration/orchastration.log` (falls back to `~/.cache/orchastration/orchastration.log`)
- Windows: `%LocalAppData%\\orchastration\\orchastration.log`

## Permissions
Runs as a normal user. It only needs read access to files you hash and write access to your user config/cache directories.

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See LICENSE and NOTICE.
