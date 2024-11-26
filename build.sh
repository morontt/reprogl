#!/usr/bin/env bash

BUILD_VERSION=$(git describe --tags)
BUILD_COMMIT=$(git rev-parse --short HEAD)

LDFLAGS="-X 'xelbot.com/reprogl/container.BuildTime=$(date -Iseconds)'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.Version=$BUILD_VERSION'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.GitRevision=$BUILD_COMMIT'"

go fmt ./...
env GOOS=linux GOARCH=amd64 go build -v -ldflags "$LDFLAGS" -tags dev ./cmd/reprogl
