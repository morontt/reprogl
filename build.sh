#!/usr/bin/env bash

BUILD_VERSION=$(git describe --tags)
BUILD_COMMIT=$(git rev-parse --short HEAD)

LDFLAGS="-X 'xelbot.com/reprogl/container.BuildTime=$(date -u +"%a, %d %b %Y %H:%M:%S %Z")'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.Version=$BUILD_VERSION'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.GitRevision=$BUILD_COMMIT'"

go fmt ./...
go build -v -ldflags "$LDFLAGS -linkmode 'external' -extldflags '-static'" -tags dev
