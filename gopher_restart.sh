#!/usr/bin/env bash

docker compose stop gopher

. ./build.sh

docker compose up -d gopher
