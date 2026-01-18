# V2 Multi-Agent Orchestration Design

## Definition and Benefits
Multi-agent orchestration coordinates multiple specialized agents in a structured workflow so complex
tasks are completed with clear role ownership and predictable handoffs. An orchestration engine acts as
the conductor, ensuring each agent runs in the right order (or in parallel when safe), with outputs flowing
between stages. This approach mirrors a well-run team, improving modularity, efficiency, and output
quality compared to a single monolithic agent. [1] [2]

Key benefits for Orchastration v2:
- Role specialization improves focus and quality of each stage (plan, build, review, doc). [1]
- The orchestration engine provides deterministic, repeatable workflows with clear traceability. [1] [2]
- Parallel execution enables faster turnaround when tasks are independent. [1]
- Shared context and run records improve auditability and handoff clarity across agents. [1]

## Core Agent Roles
Each agent implements a common interface for uniform orchestration. Initial core roles:
- PlannerAgent: Decomposes a high-level goal into a structured task plan and stores it in shared context.
  This creates an actionable roadmap for downstream stages. [1] [2]
- BuilderAgent: Executes the plan to produce concrete outputs (code, configs, artifacts), integrating with
  existing build logic where appropriate. [1] [2]
- ReviewerAgent: Validates outputs via tests, checks, or acceptance criteria and records feedback or
  status in context. This is the quality gate. [1] [2]
- DocAgent: Summarizes results and updates documentation (README, docs/ pages, usage notes), ensuring
  traceability of what was built and why. [1] [2]

## Orchestration Engine (The Conductor)
The Orchestration Engine is a centralized coordinator that:
- Accepts an orchestration definition (ordered agent list, optional parallel groups). [1] [2]
- Creates a shared OrchContext for the run.
- Invokes agents, passes shared context, captures outputs, and handles errors consistently. [1]
- Persists run metadata and key outputs in the state directory for auditability. [1]

This separation keeps agents focused on their specialty while the engine manages timing, ordering,
concurrency, and data flow. [1]

## Orchestration Patterns
Sequential pattern:
- Default pipeline: Planner -> Builder -> Reviewer -> Doc.
- Suits workflows with strict dependencies, enabling progressive refinement and predictable results. [1] [2]

Concurrent pattern:
- Used when tasks are independent and can safely run in parallel.
- The engine can fan out multiple agents and fan in their results for downstream steps. [1]
- Requires thread-safe shared context and careful aggregation of outputs.

## Shared Context Design and Persistence
Shared context (OrchContext) is a run-scoped, in-memory store used by all agents. [1] [2]
Design goals:
- Simple Set/Get API with clear, stable keys (e.g., plan.tasks, build.outputs, review.report).
- Thread-safe implementation to support concurrent runs in later steps.
- Persist important keys and run metadata to the state directory, aligned with existing task run records,
  for traceability and post-run inspection. [1]

## DevOps Workflow
We will follow an issue-driven workflow for each v2 feature:
- Create a GitHub issue per task with descriptive title and acceptance criteria. [2]
- Implement each issue on a dedicated feature branch.
- Use small, meaningful commits in imperative mood, referencing the issue number. [2]
- Merge branches sequentially into master to keep history linear and auditable. [2]
- Maintain backward compatibility with v1 commands (plan, build, doc) and cross-platform behavior.

## References
[1] Designing a Multi-Agent Orchestration Layer for Orchastration CLI.pdf
[2] Orchastration v2.0.0: Multi-Agent Orchestration Implementation.pdf
