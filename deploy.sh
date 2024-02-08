#!/bin/bash
set -e

# git add .
# git commit -m "updating changes to v0.1.3"
# git push

git tag v0.1.6
git push --tags
GOPROXY=proxy.golang.org go list -m github.com/raywall/aws-lowcode-lambda-go@v0.1.6