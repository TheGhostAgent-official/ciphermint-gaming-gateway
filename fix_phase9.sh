#!/usr/bin/env bash
set -euo pipefail

cd /workspaces/ciphermint-gaming-gateway

echo "Patching API handlers to match store signatures..."

# 1) EarnHandler: send integrationID + source to UpdateBalance
sed -i \
  's/h\.store\.UpdateBalance(r\.Context(), playerID, req\.Token, req\.Amount)/h.store.UpdateBalance(r.Context(), integrationID, playerID, req.Token, req.Source, req.Amount)/' \
  internal/api/http.go

# 2) GetPlayerHandler: send integrationID into GetPlayer
sed -i \
  's/h\.store\.GetPlayer(r\.Context(), playerID)/h.store.GetPlayer(r.Context(), integrationID, playerID)/' \
  internal/api/http.go

# 3) Format and make sure everything builds cleanly
go fmt ./...
go mod tidy
go build ./...

echo
echo "✅ Phase 9 patch complete – API handlers now match sqlstore.Store."
