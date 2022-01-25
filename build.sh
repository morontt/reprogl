#!/usr/bin/env bash

BUILD_VERSION=$(git describe --tags)

LDFLAGS="-X 'xelbot.com/reprogl/container.BuildTime=$(date -u +"%a, %d %b %Y %H:%M:%S %Z")'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.Version=$BUILD_VERSION'"

go fmt ./...
go build -v -ldflags "$LDFLAGS"
