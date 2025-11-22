#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://127.0.0.1:8080"

# Use the environment variable you exported in the shell
HDR_API_KEY="X-RACKDOG-API-KEY: ${RACKDOG_API_KEY:-}"

echo "==================================="
echo " CipherMint Gaming Gateway – RACKDOG DEMO (API Key)"
echo "==================================="
echo

echo "1) Health check..."
curl -sS "${BASE_URL}/health"
echo
echo

echo "2) Register integration 'ghostops_cod'..."
curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "${HDR_API_KEY}" \
  "${BASE_URL}/v1/game" \
  -d '{
        "id": "ghostops_cod",
        "name": "Ghost Ops -- CoD Integration",
        "company_id": ""
      }'
echo
echo

echo "3) Create player 'player123'..."
curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "${HDR_API_KEY}" \
  "${BASE_URL}/v1/game/ghostops_cod/player" \
  -d '{
        "player_id": "player123",
        "alias": "GhostPlayer"
      }'
echo
echo

echo "4) Give signup bonus in RACKDOG..."
curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "${HDR_API_KEY}" \
  "${BASE_URL}/v1/game/ghostops_cod/player/player123/earn" \
  -d '{
        "token": "RACKDOG",
        "amount": 100,
        "source": "signup_bonus"
      }'
echo
echo

echo "5) Fetch player state with balances..."
curl -sS -X GET \
  -H "${HDR_API_KEY}" \
  "${BASE_URL}/v1/game/ghostops_cod/player/player123"
echo
echo "Done."
