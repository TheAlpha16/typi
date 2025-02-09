#!/bin/bash

set -e
set -o pipefail
set -u

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
ENV_FILE="$SCRIPT_DIR/../../../api/.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "[-] env file $ENV_FILE missing."
    exit 1
fi

if docker ps -q --filter "name=typi-api" | grep -q .; then
    echo "[+] stopping and removing existing container..."
    docker rm -f typi-api
fi

if docker images -q typi-api | grep -q .; then
    echo "[+] removing old image..."
    docker rmi -f typi-api
fi

cd "$SCRIPT_DIR/../../../api" || { echo "[-] failed to cd into API directory"; exit 1; }

echo "[+] building API image..."
docker build -t typi-api .

echo "[+] running new API container..."
docker run -d --name typi-api \
    --network=host \
    --env-file "$ENV_FILE" \
    typi-api
