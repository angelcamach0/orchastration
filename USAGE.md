# Usage Guide

This guide walks through running Orchastration locally on a single machine.

## Build

Linux:
```bash
./scripts/build-linux.sh
```

Windows (PowerShell):
```powershell
.\scripts\build-windows.ps1
```

## Quick Start

1. Copy the example config to your OS config directory or pass `--config`.

Linux example:
```bash
mkdir -p ~/.config/orchastration
cp configs/config.example.toml ~/.config/orchastration/config.toml
```

Windows example (PowerShell):
```powershell
New-Item -ItemType Directory -Force -Path "$env:AppData\orchastration" | Out-Null
Copy-Item -Force configs\config.example.toml "$env:AppData\orchastration\config.toml"
```

2. List configured jobs:
```bash
./dist/orchastration list
```

3. Run a job:
```bash
./dist/orchastration run sample
```

4. Check last run status:
```bash
./dist/orchastration status
```

5. Plan and build a task:
```bash
./dist/orchastration plan list
./dist/orchastration plan create sample_task
./dist/orchastration plan status sample_task
./dist/orchastration build run sample_task
./dist/orchastration doc generate sample_task
./dist/orchastration git issue create sample_task
./dist/orchastration git branch create sample_task
```

## Commands

- `orchastration list`: show configured jobs
- `orchastration run <job-name>`: execute a job by name
- `orchastration status`: show last recorded run for each job
- `orchastration plan list`: list configured tasks
- `orchastration plan create <task>`: initialize task state
- `orchastration plan status <task>`: show task state
- `orchastration build run <task>`: execute task command
- `orchastration doc generate <task>`: generate task documentation
- `orchastration git issue create <task>`: create a GitHub issue using `gh`
- `orchastration git branch create <task>`: create a git branch for the task
- `orchastration hash --file <path>`: compute file hash
- `orchastration --help`: show help
- `orchastration --version`: show version

## Global Flags

- `--config <path>`: override config location
- `--state-dir <path>`: override state directory location

## Job Configuration

Jobs are defined in `config.toml` using argv arrays (no shell parsing):
```toml
[jobs.sample]
description = "List current directory"
command = ["ls", "-la"]
working_dir = "."
timeout_seconds = 10
env = { SAMPLE_ENV = "true" }
```

## State Records

Each run writes JSON records under the state directory:
```
state/runs/<job-name>/<timestamp>.json
state/runs/<job-name>/last.json
```

State directory locations (defaults):
- Linux: `$XDG_CACHE_HOME/orchastration/state` (falls back to `~/.cache/orchastration/state`)
- Windows: `%LocalAppData%\orchastration\state`

## Logs

Logs are JSON to stdout and a log file:
- Linux: `$XDG_CACHE_HOME/orchastration/orchastration.log` (falls back to `~/.cache/orchastration/orchastration.log`)
- Windows: `%LocalAppData%\orchastration\orchastration.log`
