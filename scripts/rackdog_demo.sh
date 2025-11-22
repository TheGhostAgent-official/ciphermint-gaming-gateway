#!/usr/bin/env bash
set -euo pipefail

echo "========================================="
echo " CipherMint Gaming Gateway – RACKDOG DEMO"
echo "========================================="

BASE_URL="http://localhost:8080"

echo
echo "1) Health check..."
curl -s "$BASE_URL/v1/health"
echo

echo
echo "2) Register integration 'ghostops_cod'..."
curl -s -X POST \
  -H "X-RACKDOG-API-KEY: $RACKDOG_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
      "id":"ghostops_cod",
      "name":"Ghost Ops -- CoD Integration",
      "company_id":""
  }' \
  "$BASE_URL/v1/game"
echo

echo
echo "3) Create player 'player123'..."
curl -s -X POST \
  -H "X-RACKDOG-API-KEY: $RACKDOG_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
      "player_id":"player123",
      "alias":"GhostPlayer"
  }' \
  "$BASE_URL/v1/game/ghostops_cod/player"
echo

echo
echo "4) Give signup bonus in RACKDOG..."
curl -s -X POST \
  -H "X-RACKDOG-API-KEY: $RACKDOG_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
      "token":"RACKDOG",
      "amount":100,
      "source":"signup_bonus"
  }' \
  "$BASE_URL/v1/game/ghostops_cod/player/player123/earn"
echo

echo
echo "5) Fetch player state with balances..."
curl -s "$BASE_URL/v1/game/ghostops_cod/player/player123"
echo

echo
echo "Done."