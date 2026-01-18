# Configuration

## Locations
- Linux: `$XDG_CONFIG_HOME/orchastration/config.toml` (falls back to `~/.config/orchastration/config.toml`)
- Windows: `%AppData%\orchastration\config.toml`

## Format
TOML

## Example
```toml
[logging]
level = "info"

[hash]
algorithm = "sha256"

[jobs.sample]
description = "List current directory"
command = ["ls", "-la"]
working_dir = "."
timeout_seconds = 10
env = { SAMPLE_ENV = "true" }

[tasks.sample_task]
description = "Example task definition"
repo = "orchastration"
working_dir = "/absolute/path"
command = ["echo", "hello"]
outputs = ["dist/example.txt"]
documents = ["README.md"]
status = "planned"

[agents.PlannerAgent]

[agents.BuilderAgent]

[agents.ReviewerAgent]

[agents.DocAgent]

[orchestrations.hello_multi_agent]
agents = ["PlannerAgent", "BuilderAgent", "ReviewerAgent", "DocAgent"]
description = "Example orchestration: plan, build, review, document"

[orchestrations.parallel_example]
steps = [
  ["PlannerAgent"],
  ["BuilderAgent", "ReviewerAgent"],
  ["DocAgent"],
]
description = "Parallel build/review step example"
```

## Options
- `logging.level`: `debug`, `info`, `warn`, `error`
- `hash.algorithm`: `sha256`, `sha1`, `sha512`
- `jobs.<name>.description`: short human description
- `jobs.<name>.command`: array form of the command and arguments (argv)
- `jobs.<name>.working_dir`: working directory for the command
- `jobs.<name>.timeout_seconds`: timeout in seconds (0 means no timeout)
- `jobs.<name>.env`: map of environment variables to add or override
- `tasks.<task>.description`: task purpose
- `tasks.<task>.repo`: `orchastration` or `external`
- `tasks.<task>.working_dir`: absolute working directory for the task
- `tasks.<task>.command`: array form of the command and arguments (argv)
- `tasks.<task>.outputs`: relative paths expected from the task
- `tasks.<task>.documents`: documentation files tied to the task
- `tasks.<task>.status`: `planned`, `in_progress`, `done`
- `agents.<name>`: reserved for agent-specific config
- `orchestrations.<name>.agents`: ordered list of agent names to run
- `orchestrations.<name>.steps`: nested agent lists (each inner list runs in parallel)
- `orchestrations.<name>.description`: human description of the orchestration

Task state is stored under `state/tasks/<task>.json` in the OS-appropriate state directory.
