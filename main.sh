#!/bin/bash

echo "🚀 Building Go App for Linux..."

APP_NAME="kasir-api"
OUTPUT_DIR="dist"

# Bersihkan output lama
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# Build Linux binary dari cmd/server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -o $OUTPUT_DIR/$APP_NAME \
  ./cmd/server

if [ $? -eq 0 ]; then
  echo "✅ Build sukses!"
  echo "📦 Output: $OUTPUT_DIR/$APP_NAME"
else
  echo "❌ Build gagal"
  exit 1
fi
