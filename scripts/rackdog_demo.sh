#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080"

# Ensure we have some API key value (dev-safe default)
: "${RACKDOG_API_KEY:=DEV_MODE}"

echo "========================================"
echo " CipherMint Gaming Gateway – RACKDOG DEMO (API Key)"
echo "========================================"
echo

# 1) Health check (health does NOT require API key)
echo "1) Health check..."
curl -s "${BASE_URL}/health"
echo
echo

# 2) Register integration 'ghostops_cod'
echo "2) Register integration 'ghostops_cod'..."
curl -s -X POST "${BASE_URL}/v1/game" \
  -H "Content-Type: application/json" \
  -H "X-RACKDOG-API-KEY: ${RACKDOG_API_KEY}" \
  -d '{
    "id": "ghostops_cod",
    "name": "Ghost Ops -- CoD Integration",
    "company_id": ""
  }'
echo
echo

# 3) Create player 'player123'
echo "3) Create player 'player123'..."
curl -s -X POST "${BASE_URL}/v1/game/ghostops_cod/player" \
  -H "Content-Type: application/json" \
  -H "X-RACKDOG-API-KEY: ${RACKDOG_API_KEY}" \
  -d '{
    "player_id": "player123",
    "alias": "GhostPlayer"
  }'
echo
echo

# 4) Give signup bonus in RACKDOG
echo "4) Give signup bonus in RACKDOG..."
curl -s -X POST "${BASE_URL}/v1/game/ghostops_cod/player/player123/earn" \
  -H "Content-Type: application/json" \
  -H "X-RACKDOG-API-KEY: ${RACKDOG_API_KEY}" \
  -d '{
    "token": "RACKDOG",
    "amount": 100,
    "source": "signup_bonus"
  }'
echo
echo

# 5) Fetch player state with balances
echo "5) Fetch player state with balances..."
curl -s "${BASE_URL}/v1/game/ghostops_cod/player/player123" \
  -H "X-RACKDOG-API-KEY: ${RACKDOG_API_KEY}"
echo
echo

echo "Done."