#!/usr/bin/env just --justfile
# To be the absolute go simp, use https://taskfile.dev (just kidding, or...)
# Disable CGo because I just want to build with pure Go

# TODO: This should be parsed from a git tag

version := "v0.0.1"
appName := "2fa"

default:
    @just --list


build:
    @go build --ldflags="-X main.Version={{ version }} -X main.AppName={{ appName }}" -o ./dist/{{ appName }} .


test *FLAGS:
    @go test ./... {{ FLAGS }}
    @go test ./otp/... {{ FLAGS }}

clean:
    rm -ri ./dist

env:
    @go env

fmt:
    @gofmt -s -w -l .

tidy:
    @cd ./otp && go mod tidy
    @go mod tidy
    @go work sync

serve-charm:
    @charm serve
