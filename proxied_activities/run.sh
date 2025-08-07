#!/bin/bash
set -e

# Kill existing processes on ports 9991 and 9992
lsof -ti:9991 | xargs kill -9 2>/dev/null || true
lsof -ti:9992 | xargs kill -9 2>/dev/null || true

# Start server in background
echo "Starting server..."
go run server/main.go &
SERVER_PID=$!

# Wait a moment for server to start
sleep 1

# Start proxy in background
echo "Starting proxy..."
go run proxy/main.go &
PROXY_PID=$!

# Wait a moment for proxy to start
sleep 1

# Start 3 workers in background
echo "Starting workers..."
go run worker/main.go >/dev/null 2>&1 &
WORKER1_PID=$!

go run worker/main.go >/dev/null 2>&1 &
WORKER2_PID=$!

go run worker/main.go >/dev/null 2>&1 &
WORKER3_PID=$!

# Wait a moment for workers to start
sleep 1

# Run starter
echo "Running starter..."
go run starter/main.go

# Clean up background processes
echo "Cleaning up..."
kill $SERVER_PID $PROXY_PID $WORKER1_PID $WORKER2_PID $WORKER3_PID 2>/dev/null
