#!/usr/bin/env bash
set -euo pipefail

cd /workspaces/ciphermint-gaming-gateway

echo "=============================================="
echo " CipherMint Gaming Gateway – RACKDOG DEMO (API Key)"
echo "=============================================="

# 1) Stop any existing gateway already on 8080
PID_8080="$(lsof -ti :8080 || true)"
if [[ -n "$PID_8080" ]]; then
  echo "Killing existing gateway on 8080 (PID $PID_8080)..."
  kill -9 "$PID_8080"
fi

# 2) Backup and remove any old SQLite DB (schema without integration_id)
if [[ -f ciphermint_gateway.db ]]; then
  TS="$(date +%s)"
  echo "Backing up old DB to ciphermint_gateway.db.bak_${TS}"
  mv ciphermint_gateway.db "ciphermint_gateway.db.bak_${TS}"
fi

# 3) Rebuild and start the gateway (fresh DB with correct schema)
echo "Running go mod tidy + go build..."
go mod tidy
go build ./...

echo "Starting gateway on :8080..."
go run main.go > gateway_api_demo.log 2>&1 &
GATEWAY_PID=$!

echo "Gateway PID: ${GATEWAY_PID}"
echo "Giving the gateway a moment to boot..."
sleep 3

# 4) Run the RackDog demo WITH the API key header
cd scripts
export RACKDOG_API_KEY="SUPERSECRET_RACKDOG_KEY_001"

echo "Running RackDog API-key demo script..."
./rackdog_demo.sh

cd ..

# 5) Cleanly stop the background gateway
echo "Stopping gateway PID ${GATEWAY_PID}..."
kill -9 "${GATEWAY_PID}" || true

echo "=============================================="
echo " API-key RackDog demo complete."
echo " - Logs: gateway_api_demo.log"
echo "=============================================="