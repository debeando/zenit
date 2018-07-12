#!/usr/bin/env bash

set -e

go test -v ./...
go build -ldflags "-s -w" -o zenit main.go
