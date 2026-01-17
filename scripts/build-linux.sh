#!/usr/bin/env bash
set -euo pipefail

APP_NAME=orchastration
VERSION=${VERSION:-dev}
COMMIT=${COMMIT:-unknown}
BUILDTIME=${BUILDTIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

mkdir -p dist

go build -trimpath \
  -ldflags "-X 'orchastration/internal/version.Version=${VERSION}' -X 'orchastration/internal/version.Commit=${COMMIT}' -X 'orchastration/internal/version.BuildTime=${BUILDTIME}'" \
  -o dist/${APP_NAME} \
  ./cmd/${APP_NAME}
