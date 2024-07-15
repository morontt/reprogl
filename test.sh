#!/usr/bin/env bash

LDFLAGS="-X 'xelbot.com/reprogl/container.BuildTime=$(date -Iseconds)'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.Version=9.9.9'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.GitRevision=f0f0f0'"

go test -tags=dev -ldflags="$LDFLAGS" ./...
