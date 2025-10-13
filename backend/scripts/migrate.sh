#!/bin/bash
set -euo pipefail

echo "🗄️ Running database migrations..."

# Check if DATABASE_URL is set
if [ -z "${DATABASE_URL:-}" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "✅ DATABASE_URL is set"

# Ensure curl exists (installed via nixpacks), otherwise fail fast
if ! command -v curl >/dev/null 2>&1; then
    echo "❌ ERROR: curl is not available. Ensure build environment includes curl."
    exit 1
fi

BIN_DIR="$(pwd)/bin"
MIGRATE_BIN="$BIN_DIR/migrate"

# Install prebuilt golang-migrate binary into local ./bin if missing
if [ ! -x "$MIGRATE_BIN" ]; then
    echo "📦 Downloading golang-migrate binary..."
    mkdir -p "$BIN_DIR"
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64) ARCH=amd64 ;;
        aarch64|arm64) ARCH=arm64 ;;
        *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    VERSION=v4.19.0
    URL="https://github.com/golang-migrate/migrate/releases/download/${VERSION}/migrate.${OS}-${ARCH}.tar.gz"
    TMP_DIR=$(mktemp -d)
    curl -fsSL "$URL" -o "$TMP_DIR/migrate.tar.gz"
    tar -xzf "$TMP_DIR/migrate.tar.gz" -C "$TMP_DIR"
    # The archive contains multiple files; pick the 'migrate' binary
    if [ -f "$TMP_DIR/migrate" ]; then
        mv "$TMP_DIR/migrate" "$MIGRATE_BIN"
        chmod +x "$MIGRATE_BIN"
    else
        # Some archives unpack to a folder
        BIN_PATH=$(find "$TMP_DIR" -maxdepth 2 -type f -name migrate | head -n1 || true)
        if [ -z "$BIN_PATH" ]; then
            echo "❌ Failed to locate migrate binary in the downloaded archive"
            ls -la "$TMP_DIR" || true
            exit 1
        fi
        mv "$BIN_PATH" "$MIGRATE_BIN"
        chmod +x "$MIGRATE_BIN"
    fi
    rm -rf "$TMP_DIR"
    echo "✅ golang-migrate installed to $MIGRATE_BIN"
fi

# Ensure our local bin is in PATH for any subcommands
export PATH="$BIN_DIR:$PATH"

# Run migrations
echo "🚀 Applying migrations..."
"$MIGRATE_BIN" -database "$DATABASE_URL" -path db/migrations up

echo "✅ All migrations applied successfully!"
