#!/bin/bash

set -e
set -o pipefail
set -u

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
ENV_FILE="$SCRIPT_DIR/../../../cron/.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "[-] env file $ENV_FILE missing."
    exit 1
fi

if docker ps -q --filter "name=typi-cron" | grep -q .; then
    echo "[+] stopping and removing existing container..."
    docker rm -f typi-cron
fi

if docker images -q typi-cron | grep -q .; then
    echo "[+] removing old image..."
    docker rmi -f typi-cron
fi

cd "$SCRIPT_DIR/../../../cron" || { echo "[-] failed to cd into cron directory"; exit 1; }

echo "[+] building cron image..."
docker build -t typi-cron .

echo "[+] running new cron container..."
docker run -d --name typi-cron \
    --network=host \
    --env-file "$ENV_FILE" \
    typi-cron
