#!/bin/bash
set -e

echo "🗄️ Running database migrations..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "✅ DATABASE_URL is set"

# Set GOPATH and add to PATH
export GOPATH="${GOPATH:-$HOME/go}"
export PATH="$GOPATH/bin:$PATH"

# Install golang-migrate if not present
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    echo "✅ golang-migrate installed to $GOPATH/bin/migrate"
fi

# Verify migrate is available
if ! command -v migrate &> /dev/null; then
    echo "❌ ERROR: migrate command still not found after installation"
    echo "PATH: $PATH"
    echo "GOPATH: $GOPATH"
    ls -la "$GOPATH/bin/" || echo "GOPATH/bin directory not found"
    exit 1
fi

# Run migrations
echo "🚀 Applying migrations..."
migrate -database "$DATABASE_URL" -path db/migrations up

echo "✅ All migrations applied successfully!"
