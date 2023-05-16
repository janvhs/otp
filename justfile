#!/usr/bin/env just --justfile
# To be the absolute go simp, use https://taskfile.dev (just kidding, or...)
# Disable CGo because I just want to build with pure Go

export CGO_ENABLED := "0"
export CHARM_SERVER_DATA_DIR := "./charm_data"
export CHARM_HOST := "localhost"

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

test *FLAGS:
    @go test ./... {{ FLAGS }}
    @go test ./otp/... {{ FLAGS }}

clean:
    rm -ri ./dist

env:
    @go env

fmt:
    @gofmt -s -w -l .
    @just tidy

tidy:
    @cd ./tools && go mod tidy
    @cd ./otp && go mod tidy
    @go mod tidy
    @go work sync

goose *FLAGS:
    @mkdir -p ./migrations
    @echo "for commands run just goose --help"
    go run github.com/pressly/goose/v3/cmd/goose -dir migrations/ sqlite3 {{ FLAGS }}

goose-init db *FLAGS:
    @mkdir -p ./migrations
    @echo "for commands run just goose --help"
    go run github.com/pressly/goose/v3/cmd/goose -dir migrations/ sqlite3 {{ db }} create init sql {{ FLAGS }}

serve-charm:
    @charm serve
