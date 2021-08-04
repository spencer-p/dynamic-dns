#!/usr/bin/env bash

set -o pipefail -o xtrace -o errexit

# Docker requires root on my system. I awkwardly have to propagate config to the
# root user, but use kubectl configured as my normal user.
sudo KO_DOCKER_REPO=docker.io/spencerjp $(which ko) resolve -f conf/cronjob.yaml --platform=linux/arm64 | kubectl apply -f -