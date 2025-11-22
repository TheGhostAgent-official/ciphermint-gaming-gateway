#!/bin/bash
set -euo pipefail

cd /workspaces/ciphermint-gaming-gateway

echo "==============================="
echo " CipherMint Gaming Gateway -- RACKDOG DEMO (API Key)"
echo "==============================="

# 1) Stop any existing gateway
PID_8080="$(lsof -ti :8080 || true)"
if [ -n "$PID_8080" ]; then kill -9 "$PID_8080"; fi

# 2) Remove bad DB and start fresh
if [ -f ciphermint_gateway.db ]; then
    mv ciphermint_gateway.db ciphermint_gateway.db.bak_$(date +%s)
fi

# 3) Rebuild
go mod tidy
go build ./...

# 4) Start gateway WITH the API KEY injected into the SAME process
RACKDOG_API_KEY="SUPERSECRET_RACKDOG_KEY_001" \
    go run main.go > gateway_api_demo.log 2>&1 &

GATEWAY_PID=$!
sleep 3

# 5) Run the RackDawg demo
cd scripts
./rackdog_demo.sh
cd ..

# 6) Clean shutdown
kill -9 "$GATEWAY_PID"

echo "Done."