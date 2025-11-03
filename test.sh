#!/bin/bash

# Integration test script for BSS
# This script starts the required services and runs integration tests

echo "=== BSS Integration Tests ==="
echo "Starting PostgreSQL service..."
./dev up -d postgres

echo "Waiting for PostgreSQL to be ready..."
sleep 5

echo "Running integration tests..."
echo ""

# Export environment variables for testing
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=bss
export POSTGRES_SSLMODE=disable

# Run the tests (not in short mode to include integration tests)
go test -count=1 -v ./src/...

echo ""
echo "=== Integration tests completed ==="
echo "Stopping PostgreSQL service..."
./dev down -v postgres