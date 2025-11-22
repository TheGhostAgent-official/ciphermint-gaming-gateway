#!/usr/bin/env bash
set -euo pipefail

cd /workspaces/ciphermint-gaming-gateway

echo "=== Phase 8: Build & package RackDOG gateway ==="

# A. Clean & build
go fmt ./...
go mod tidy
mkdir -p dist
go build -o dist/ciphermint-gaming-gateway main.go

# B. Prepare bundle folder
BUNDLE_DIR="dist/rackdog_gateway_bundle"
rm -rf "$BUNDLE_DIR"
mkdir -p "$BUNDLE_DIR"

# Core artifacts
cp docs/RACKDOG_GATEWAY_CONTRACT.md "$BUNDLE_DIR/"
cp scripts/rackdog_demo.sh "$BUNDLE_DIR/"

# Optional: top-level README if it exists
if [ -f README.md ]; then
  cp README.md "$BUNDLE_DIR/"
fi

# C. Create archives (tar.gz is guaranteed; zip is optional)
cd dist

# Tarball (always created)
tar -czf rackdog_gateway_bundle.tar.gz rackdog_gateway_bundle

# Zip archive if 'zip' is available
if command -v zip >/dev/null 2>&1; then
  rm -f rackdog_gateway_bundle.zip
  zip -r rackdog_gateway_bundle.zip rackdog_gateway_bundle >/dev/null
fi

echo "=== Phase 8 complete ==="
echo "Artifacts:"
echo "  dist/rackdog_gateway_bundle/            (folder)"
echo "  dist/rackdog_gateway_bundle.tar.gz      (tarball)"
echo "  dist/rackdog_gateway_bundle.zip         (zip, if zip is installed)"
