# Architecture

## Overview
Orchastration follows a small, layered layout that cleanly separates CLI entrypoints, application logic, configuration, logging, OS-specific path handling, and execution state.

## Internal Architecture (v2)

![Orchastration Internal Architecture v2](docs/diagrams/Orchastration%20F-Internal%20Architecture%20v2.svg)

This diagram maps directly to the repository layout: `cmd/` hosts the CLI entrypoint, while `internal/app` contains the command handlers and orchestration flow. The new Orchestration Engine in `internal/orchestrator` coordinates the multi-agent pipeline using the Agent Registry in `internal/agent`. State persistence lives in `internal/state`, which writes task records and orchestration run records under the OS-native state directory. Support modules like `internal/config`, `internal/logging`, and `internal/platform` provide configuration parsing, structured logging, and OS-aware paths without leaking those concerns into the command logic. The result is a clear separation of concerns that keeps orchestration behavior stable and testable.

## Multi-Agent Workflow

![Multi-Agent Workflow](docs/diagrams/Orchastration%20E-Multi-Agent%20Workflow.svg)

The orchestration engine runs Planner -> Builder -> Reviewer -> Doc by default, sharing a thread-safe OrchContext that carries plan, outputs, review results, and documentation paths across stages.

## Context Sharing (Parallel)

![Context Sharing in Parallel](docs/diagrams/Orchastration%20G-Context%20Sharing%20%28Parallel%29.svg)

Parallel agent groups can run concurrently when defined in config `steps`. Each agent reads and writes shared context safely, and results are persisted with the orchestration run record.

### External Repository Interaction

![Orchastration External Repository Interaction](docs/diagrams/Orchastration%20H-External%20Repository%20Interaction.svg)

This diagram shows Orchastration executing tasks inside external repositories without owning their domain logic. Task definitions provide an absolute `working_dir` and argv-style `command`, and Orchastration runs them as-is while capturing logs and artifacts. Outputs and documentation are stored by the external repo, while Orchastration only records metadata and status in its own state directory. This keeps the orchestration agent decoupled from the repository’s code and build systems.

## Layers
- `cmd/orchastration`: CLI entrypoint, delegates to application layer.
- `internal/app`: Command parsing and orchestration for each CLI command (jobs, tasks, agents, orchestrations).
- `internal/agent`: Agent interface, registry, and core agent implementations.
- `internal/config`: Config structs and TOML loading.
- `internal/logging`: Structured logging setup.
- `internal/orchestrator`: Orchestration engine coordinating agent runs.
- `internal/platform`: OS-aware config and log paths.
- `internal/state`: Execution record persistence.
- `internal/taskflow`: Shared task planning/building/documentation logic used by CLI and agents.
- `internal/version`: Build-time version metadata.

## Execution Flow
1. CLI parses global flags and dispatches to a command.
2. Config is loaded from the OS-appropriate location (or a user-provided path).
3. Logging is configured for stdout and file output.
4. Commands execute with structured logs and explicit error handling.
5. Orchestration runs create a shared OrchContext and persist run records under the OS-appropriate state directory.

## Extensibility
Add new commands by creating a new `runX` function in `internal/app` and wiring it in the command switch. If a command needs OS-specific behavior, add a helper to `internal/platform`.
