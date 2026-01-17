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
```

## Options
- `logging.level`: `debug`, `info`, `warn`, `error`
- `hash.algorithm`: `sha256`, `sha1`, `sha512`
- `jobs.<name>.description`: short human description
- `jobs.<name>.command`: array form of the command and arguments (argv)
- `jobs.<name>.working_dir`: working directory for the command
- `jobs.<name>.timeout_seconds`: timeout in seconds (0 means no timeout)
- `jobs.<name>.env`: map of environment variables to add or override
