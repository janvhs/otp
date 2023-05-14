#!/usr/bin/env just --justfile
# To be the absolute go simp, use https://taskfile.dev (just kidding, or...)

# Disable CGo because I just want to build with pure Go
export CGO_ENABLED := "0"

# TODO: This should be parsed from a git tag
version := "v0.0.1"

appName := "2fa"

default:
  @just --list

run *FLAGS:
    @go run . {{ FLAGS }}

build:
    @go build --ldflags="-X main.Version={{ version }} -X main.AppName={{ appName }}" -o ./dist/{{ appName }} .

build-run *FLAGS:
    @just build
    @./dist/{{ appName }} {{ FLAGS }}

clean:
    rm -ri ./dist

env:
    @go env

fmt:
    @go fmt .
    @go mod tidy

goose *FLAGS:
    @mkdir -p ./migrations
    @echo "for commands run just goose --help"
    go run github.com/pressly/goose/v3/cmd/goose -dir migrations/ sqlite3 {{ FLAGS }}

goose-init db *FLAGS:
    @mkdir -p ./migrations
    @echo "for commands run just goose --help"
    go run github.com/pressly/goose/v3/cmd/goose -dir migrations/ sqlite3 {{ db }} create init sql {{ FLAGS }}