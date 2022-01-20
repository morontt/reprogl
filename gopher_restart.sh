#!/bin/bash

docker-compose stop source

go fmt ./...
go build

docker-compose up -d source
