#!/usr/bin/env bash

BUILD_VERSION=$(git describe --tags)
BUILD_COMMIT=$(git rev-parse --short HEAD)
SRI_JS=$(cat public/assets/js/reprogl.min.js | openssl dgst -sha256 -binary | openssl base64 -A)
SRI_CSS=$(cat public/assets/css/reprogl.min.css | openssl dgst -sha256 -binary | openssl base64 -A)

LDFLAGS="-X 'xelbot.com/reprogl/container.BuildTime=$(date -u +"%a, %d %b %Y %H:%M:%S %Z")'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.Version=$BUILD_VERSION'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.GitRevision=$BUILD_COMMIT'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.SubresourceIntegrityJS=$SRI_JS'"
LDFLAGS="$LDFLAGS -X 'xelbot.com/reprogl/container.SubresourceIntegrityCSS=$SRI_CSS'"

go fmt ./...
go build -v -ldflags "$LDFLAGS" -tags dev
