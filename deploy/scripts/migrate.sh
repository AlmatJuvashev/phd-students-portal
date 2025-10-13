#!/bin/bash
set -e

echo "🚀 Running database migrations..."

cd backend

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "✅ DATABASE_URL is set"

# Run migrations using psql
echo "📦 Applying migrations..."

# Apply each migration file
for file in db/migrations/*.up.sql; do
    echo "  → Running $(basename $file)..."
    psql "$DATABASE_URL" -f "$file"
done

echo "✅ All migrations applied successfully!"
