#!/bin/bash
set -e

sam build && sam local start-api --skip-pull-image --warm-containers eager --docker-network sam-local-network --debug