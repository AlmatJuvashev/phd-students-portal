#!/bin/bash

# Safe Kill Ports Script
# Usage: ./kill_ports.sh <port1> <port2> ...

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 <port1> <port2> ..."
    exit 1
fi

PORTS=("$@")

# Define forbidden process names (case-insensitive grep)
FORBIDDEN_NAMES="Antigravity|Code|VSCode|Electron|Google Chrome"

echo "Checking ports: ${PORTS[*]}"

for port in "${PORTS[@]}"; do
    # 1. Check for Docker containers on this port
    # We use --filter publish=port to find containers mapping this port
    container_id=$(docker ps --filter "publish=$port" --format "{{.ID}}")
    
    if [ ! -z "$container_id" ]; then
        echo "🐳 Found Docker container on port $port: $container_id"
        echo "   Stopping container..."
        docker rm -f "$container_id" > /dev/null
        echo "   ✅ Container stopped."
        # If we found a container, we assume it was the primary holder. 
        # But we still check for local processes just in case (e.g. if docker command failed or didn't clear it).
    fi

    # 2. Check for local processes
    # -t: terse (PID only)
    # -i: internet address
    pids=$(lsof -ti:$port 2>/dev/null)
    
    if [ -z "$pids" ]; then
        continue
    fi
    
    for pid in $pids; do
        # Get full command line for check
        command=$(ps -p $pid -o command=)
        
        # Check if it matches forbidden names OR is the Docker backend itself
        # (Docker backend often shows up in lsof but we shouldn't kill it directly, we handled specific containers above)
        if echo "$command" | grep -Eiq "$FORBIDDEN_NAMES|com.docker.backend|Docker"; then
            echo "⚠️  Skipping protected/system process on port $port: PID $pid ($command)"
        else
            echo "✅ Killing process on port $port: PID $pid ($command)"
            kill -9 $pid 2>/dev/null || true
        fi
    done
done

echo "Done clearing ports."
