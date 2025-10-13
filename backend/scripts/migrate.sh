#!/bin/bash
set -e

echo "🗄️ Running database migrations..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "✅ DATABASE_URL is set"

# Install golang-migrate if not present
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Run migrations
echo "🚀 Applying migrations..."
cd /app/backend
migrate -database "$DATABASE_URL" -path db/migrations up

echo "✅ All migrations applied successfully!"
