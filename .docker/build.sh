#!/bin/bash

cd /local/workspace/aws-lowcode-lambda-go
go build -o /var/runtime/bootstrap cmd/lambda/main.go
# /var/task/bootstrap