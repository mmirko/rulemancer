#!/usr/bin/env bash

# Kill the rulemancer server using the PID stored in RULEMANCER_PID

# Check if RULEMANCER_PID is set
if [ -z "$RULEMANCER_PID" ]; then
    echo "Error: RULEMANCER_PID environment variable not set"
    echo "The server may not have been started with launch.sh, or the variable is not exported"
    return 1 2>/dev/null || exit 1
fi

# Check if the process exists
if ! ps -p "$RULEMANCER_PID" > /dev/null 2>&1; then
    echo "Warning: No process found with PID $RULEMANCER_PID"
    unset RULEMANCER_PID
    unset API_TOKEN
    return 1 2>/dev/null || exit 1
fi

echo "Killing rulemancer server with PID: $RULEMANCER_PID"

# Try graceful shutdown first (SIGTERM)
kill -TERM "$RULEMANCER_PID" 2>/dev/null

# Wait up to 5 seconds for graceful shutdown
for i in {1..10}; do
    if ! ps -p "$RULEMANCER_PID" > /dev/null 2>&1; then
        echo "Server stopped successfully"
        unset RULEMANCER_PID
        unset API_TOKEN
        return 0 2>/dev/null || exit 0
    fi
    sleep 0.5
done

# Force kill if still running
echo "Server did not stop gracefully, forcing shutdown..."
kill -KILL "$RULEMANCER_PID" 2>/dev/null

# Verify it's dead
if ! ps -p "$RULEMANCER_PID" > /dev/null 2>&1; then
    echo "Server forcefully stopped"
    unset RULEMANCER_PID
    unset API_TOKEN
    return 0 2>/dev/null || exit 0
else
    echo "Error: Failed to stop server"
    return 1 2>/dev/null || exit 1
fi
