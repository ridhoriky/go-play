#!/bin/bash
# ─────────────────────────────────────────────────────────────────────────────
# generate-jwt-keys.sh
# Generate ECDSA ES256 key pairs untuk JWT access & refresh token.
# Output: 4 file PEM di direktori ./keys/ (gitignore-d)
# Usage: bash scripts/generate-jwt-keys.sh
# ─────────────────────────────────────────────────────────────────────────────

set -euo pipefail

KEYS_DIR="./keys"
mkdir -p "$KEYS_DIR"

echo "🔑 Generating ECDSA P-256 key pairs for JWT (ES256)..."
echo ""

# ── Access Token Keys ──────────────────────────────────────────────────────
echo "  [1/4] Generating access token private key..."
openssl ecparam -genkey -name prime256v1 -noout \
  | openssl pkcs8 -topk8 -nocrypt \
  > "$KEYS_DIR/access_private.pem"

echo "  [2/4] Extracting access token public key..."
openssl ec -in "$KEYS_DIR/access_private.pem" -pubout \
  > "$KEYS_DIR/access_public.pem" 2>/dev/null

# ── Refresh Token Keys ─────────────────────────────────────────────────────
echo "  [3/4] Generating refresh token private key..."
openssl ecparam -genkey -name prime256v1 -noout \
  | openssl pkcs8 -topk8 -nocrypt \
  > "$KEYS_DIR/refresh_private.pem"

echo "  [4/4] Extracting refresh token public key..."
openssl ec -in "$KEYS_DIR/refresh_private.pem" -pubout \
  > "$KEYS_DIR/refresh_public.pem" 2>/dev/null

echo ""
echo "✅ Keys generated:"
echo "   $KEYS_DIR/access_private.pem"
echo "   $KEYS_DIR/access_public.pem"
echo "   $KEYS_DIR/refresh_private.pem"
echo "   $KEYS_DIR/refresh_public.pem"
echo ""
echo "─────────────────────────────────────────────────────────────────"
echo "📋 Add these lines directly to your .env file:"
echo "─────────────────────────────────────────────────────────────────"
echo ""
printf "JWT_ACCESS_PRIVATE_KEY=\""
awk '{printf "%s\\n", $0}' "$KEYS_DIR/access_private.pem"
printf "\"\n\n"

printf "JWT_ACCESS_PUBLIC_KEY=\""
awk '{printf "%s\\n", $0}' "$KEYS_DIR/access_public.pem"
printf "\"\n\n"

printf "JWT_REFRESH_PRIVATE_KEY=\""
awk '{printf "%s\\n", $0}' "$KEYS_DIR/refresh_private.pem"
printf "\"\n\n"

printf "JWT_REFRESH_PUBLIC_KEY=\""
awk '{printf "%s\\n", $0}' "$KEYS_DIR/refresh_public.pem"
printf "\"\n"
echo ""
echo "─────────────────────────────────────────────────────────────────"
echo "⚠️  IMPORTANT: Add $KEYS_DIR/ to .gitignore!"
