#!/bin/bash
set -e

# Compila o código fonte para gerar o arquivo bootstrap
# GOOS=linux GOARCH=amd64 go build -o ../bin/bootstrap cmd/lambda/main.go

# Inicializa o container com a lambda
docker compose -f .docker/docker-compose.yaml up -d

# Inicializa o API Gateway
# sam local start-api --template .docker/template.yaml --skip-pull-image --warm-containers LAZY

# curl -X POST http://127.0.0.1:3000/ -d '{"teste":123}'