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
```

## Options
- `logging.level`: `debug`, `info`, `warn`, `error`
- `hash.algorithm`: `sha256`, `sha1`, `sha512`
