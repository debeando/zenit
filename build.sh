#!/usr/bin/env bash

set -e

go test -cover ./...
go build -ldflags "-s -w" -o zenit main.go
