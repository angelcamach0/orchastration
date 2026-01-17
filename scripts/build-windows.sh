#!/usr/bin/env bash
set -euo pipefail

APP_NAME=orchastration
VERSION=${VERSION:-dev}
COMMIT=${COMMIT:-unknown}
BUILDTIME=${BUILDTIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

mkdir -p dist

GOOS=windows GOARCH=amd64 go build -trimpath \
  -ldflags "-X 'orchastration/internal/version.Version=${VERSION}' -X 'orchastration/internal/version.Commit=${COMMIT}' -X 'orchastration/internal/version.BuildTime=${BUILDTIME}'" \
  -o dist/${APP_NAME}.exe \
  ./cmd/${APP_NAME}
