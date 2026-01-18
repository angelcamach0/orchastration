# Architecture

## Overview
Orchastration follows a small, layered layout that cleanly separates CLI entrypoints, application logic, configuration, logging, OS-specific path handling, and execution state.

## Internal Architecture

![Orchastration Internal Architecture](docs/diagrams/Orchastration%20B-Internal%20Architecture.svg)

This diagram maps directly to the repository layout: `cmd/` hosts the CLI entrypoint, while `internal/app` contains the command handlers and orchestration flow. State persistence lives in `internal/state`, which writes task records and run records under the OS-native state directory. Support modules like `internal/config`, `internal/logging`, and `internal/platform` provide configuration parsing, structured logging, and OS-aware paths without leaking those concerns into the command logic. The result is a clear separation of concerns that keeps orchestration behavior stable and testable.

### External Repository Interaction

![Orchastration BioZero Interaction](docs/diagrams/Orchastration%20D-Orchastration%20%E2%86%94%20BioZero%20Interaction.svg)

This diagram shows Orchastration executing tasks inside external repositories without owning their domain logic. Task definitions provide an absolute `working_dir` and argv-style `command`, and Orchastration runs them as-is while capturing logs and artifacts. Outputs and documentation are stored by the external repo, while Orchastration only records metadata and status in its own state directory. This keeps the orchestration agent decoupled from the repository’s code and build systems.

## Layers
- `cmd/orchastration`: CLI entrypoint, delegates to application layer.
- `internal/app`: Command parsing and orchestration for each CLI command (jobs and workflow tasks).
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
5. Job and task execution records are persisted under the OS-appropriate state directory.

## Extensibility
Add new commands by creating a new `runX` function in `internal/app` and wiring it in the command switch. If a command needs OS-specific behavior, add a helper to `internal/platform`.
