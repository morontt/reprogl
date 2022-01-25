#!/usr/bin/env bash

docker-compose stop source

. ./build.sh

docker-compose up -d source
