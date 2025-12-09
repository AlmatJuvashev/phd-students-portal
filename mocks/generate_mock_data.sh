#!/usr/bin/env bash
set -euo pipefail

# Navigate to backend directory (where go.mod is)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"

cd "$BACKEND_DIR"
go run cmd/mock/main.go
