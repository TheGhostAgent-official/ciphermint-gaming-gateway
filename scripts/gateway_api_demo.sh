#!/usr/bin/env bash
set -euo pipefail

cd /workspaces/ciphermint-gaming-gateway

# Ensure we have an API key value; if none is set in the shell, use DEV_MODE
: "${RACKDOG_API_KEY:=DEV_MODE}"

echo "==============================================="
echo " CipherMint Gaming Gateway – RACKDOG DEMO (API Key)"
echo "==============================================="

# 1) Stop any gateway already on 8080
PID_8080="$(lsof -t -i :8080 || true)"
if [ -n "${PID_8080}" ]; then
  echo "Killing existing gateway on 8080 (PID: ${PID_8080})..."
  kill -9 "${PID_8080}" || true
fi

# 2) Backup and remove any old SQLite DB (schema without integration_id)
if [ -f ciphermint_gateway.db ]; then
  echo "Backing up old DB..."
  mv ciphermint_gateway.db "ciphermint_gateway.db.bak_$(date +%s)"
fi

# 3) Rebuild and start the gateway with the API key in its environment
echo "Running go mod tidy + go build..."
go mod tidy
go build ./...

echo "Starting gateway on port 8080 with RACKDOG_API_KEY=${RACKDOG_API_KEY}..."
RACKDOG_API_KEY="${RACKDOG_API_KEY}" go run main.go > gateway_api_demo.log 2>&1 &
GATEWAY_PID=$!

# 4) Give the gateway a moment to boot
sleep 3

# 5) Run the RackDOG demo (this sends X-RACKDOG-API-KEY)
cd scripts
./rackdog_demo.sh
cd ..

# 6) Cleanly stop the background gateway
echo "Stopping gateway (PID: ${GATEWAY_PID})..."
kill -9 "${GATEWAY_PID}" || true

echo "==============================================="
echo " RACKDOG API key demo complete."
echo "==============================================="
