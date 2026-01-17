# Orchastration

Orchastration is a cross-platform CLI for small operational workflows. The initial command computes file hashes for integrity checks and automation pipelines.

## Features
- Native binaries for Windows and Linux
- Structured logging to stdout and a log file
- OS-appropriate config and log locations
- Simple, extensible command layout

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
```

## Configuration
The config file is optional. If missing, defaults are used.
- Linux: `$XDG_CONFIG_HOME/orchastration/config.toml` (falls back to `~/.config/orchastration/config.toml`)
- Windows: `%AppData%\orchastration\config.toml`

Example config: `configs/config.example.toml`

## Logging
Logs are JSON and written to stdout and a log file:
- Linux: `$XDG_CACHE_HOME/orchastration/orchastration.log` (falls back to `~/.cache/orchastration/orchastration.log`)
- Windows: `%LocalAppData%\orchastration\orchastration.log`

## Permissions
Runs as a normal user. It only needs read access to files you hash and write access to your user config/cache directories.
