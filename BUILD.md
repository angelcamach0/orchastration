# Build

## Prerequisites
- Go 1.21+

## Linux
```bash
./scripts/build-linux.sh
```
Output: `dist/orchastration`

## Windows (from Windows)
```powershell
.\scripts\build-windows.ps1
```
Output: `dist\orchastration.exe`

## Windows (cross-compile from Linux)
```bash
./scripts/build-windows.sh
```

## Version metadata
Set these environment variables before running a build script:
- `VERSION`
- `COMMIT`
- `BUILDTIME`

Example:
```bash
VERSION=2.0.0 COMMIT=$(git rev-parse --short HEAD) BUILDTIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) ./scripts/build-linux.sh
```
