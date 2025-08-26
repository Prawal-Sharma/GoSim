#!/bin/bash

# GoSim Start Script

echo "🎯 Starting GoSim - Interactive Go Learning Simulator"
echo "=================================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "Please install Go from: https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"

# Install dependencies
echo "📦 Installing dependencies..."
go mod download

# Build the server
echo "🔨 Building server..."
go build -o gosim cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting server on http://localhost:8080"
    echo "=================================================="
    echo "Press Ctrl+C to stop the server"
    echo ""
    
    # Run the server
    ./gosim
else
    echo "❌ Build failed. Please check for errors above."
    exit 1
fi