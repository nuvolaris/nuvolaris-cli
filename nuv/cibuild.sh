#!/bin/bash
TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo latest)
go build -ldflags "-X main.CLIVersion=$TAG" -tags subcommands -o nuv
