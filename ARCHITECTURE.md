# Architecture

## Overview
Orchastration follows a small, layered layout that cleanly separates CLI entrypoints, application logic, configuration, logging, OS-specific path handling, and execution state.

## Layers
- `cmd/orchastration`: CLI entrypoint, delegates to application layer.
- `internal/app`: Command parsing and orchestration for each CLI command.
- `internal/config`: Config structs and TOML loading.
- `internal/logging`: Structured logging setup.
- `internal/platform`: OS-aware config and log paths.
- `internal/state`: Execution record persistence.
- `internal/version`: Build-time version metadata.

## Execution Flow
1. CLI parses global flags and dispatches to a command.
2. Config is loaded from the OS-appropriate location (or a user-provided path).
3. Logging is configured for stdout and file output.
4. Command runs with structured logs and explicit error handling.
5. Job execution records are persisted under the OS-appropriate state directory.

## Extensibility
Add new commands by creating a new `runX` function in `internal/app` and wiring it in the command switch. If a command needs OS-specific behavior, add a helper to `internal/platform`.
