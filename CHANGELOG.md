# Changelog

## [2.0.0] - 2026-01-18
### Added
- Multi-agent orchestration engine with sequential and parallel agent group support.
- Agent interface, registry, and `orchastration agent list` CLI command.
- Orchestration config sections with `agents` and `steps`, plus orchestration list/run commands.
- Shared OrchContext with persisted orchestration run records and context snapshots.
- Planner/Builder/Reviewer/Doc agents with executable logic wired to plan/build/doc workflows.
- New diagrams and documentation for v2 architecture and orchestration flows.

### Notes
- Multi-agent orchestration enables specialized agents to collaborate on complex tasks, improving
  modularity and maintainability via clear role separation and shared context.
