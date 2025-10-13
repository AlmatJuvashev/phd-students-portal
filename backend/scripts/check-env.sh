#!/bin/bash
set -e

echo "🔍 Environment Check at Startup"
echo "================================"
echo "FRONTEND_BASE: ${FRONTEND_BASE:-NOT SET}"
echo "PORT: ${PORT:-NOT SET}"
echo "GIN_MODE: ${GIN_MODE:-NOT SET}"
echo "DATABASE_URL: ${DATABASE_URL:+***SET***}"
echo "JWT_SECRET: ${JWT_SECRET:+***SET***}"
echo "================================"
echo ""
