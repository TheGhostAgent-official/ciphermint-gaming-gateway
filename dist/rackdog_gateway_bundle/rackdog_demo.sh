#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080"
GAME_ID="ghostops_cod"
PLAYER_ID="player123"
ALIAS="GhostPlayer"
TOKEN="RACKDOG"

echo "========================================"
echo " CipherMint Gaming Gateway – RACKDOG DEMO"
echo "========================================"
echo

echo "1) Health check..."
curl -s "${BASE_URL}/health"
echo; echo

echo "2) Register integration '${GAME_ID}'..."
curl -s -X POST "${BASE_URL}/v1/game" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\":\"${GAME_ID}\",
    \"name\":\"Ghost Ops -- CoD Integration\",
    \"company_id\":\"\"
  }"
echo; echo

echo "3) Create player '${PLAYER_ID}'..."
curl -s -X POST "${BASE_URL}/v1/game/${GAME_ID}/player" \
  -H "Content-Type: application/json" \
  -d "{
    \"player_id\":\"${PLAYER_ID}\",
    \"alias\":\"${ALIAS}\"
  }"
echo; echo

echo "4) Give signup bonus in ${TOKEN}..."
curl -s -X POST "${BASE_URL}/v1/game/${GAME_ID}/player/${PLAYER_ID}/earn" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\":\"${TOKEN}\",
    \"amount\":100,
    \"source\":\"signup_bonus\"
  }"
echo; echo

echo "5) Fetch player state with balances..."
curl -s "${BASE_URL}/v1/game/${GAME_ID}/player/${PLAYER_ID}"
echo; echo

echo "Done."
