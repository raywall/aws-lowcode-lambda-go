#!/bin/bash
set -e

GOOS=linux GOARCH=amd64 go build -o bin/bootstrap cmd/lambda/main.go
sam local start-api --debug --skip-pull-image --warm-containers LAZY
# curl -X POST http://127.0.0.1:3000/ -d '{"teste":123}'