#!/bin/bash

echo "=== Testing MCP Service Compilation ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

echo "Go version: $(go version)"

# Navigate to MCP service directory
cd "$(dirname "$0")"
echo "Working directory: $(pwd)"

echo ""
echo "=== Step 1: Clean and download dependencies ==="
go clean -cache
go mod tidy

echo ""
echo "=== Step 2: Check for syntax errors ==="
go vet ./...

echo ""
echo "=== Step 3: Build test program ==="
if go build -o test_compile ./test_compile.go; then
    echo "✅ Test compilation successful!"
    ./test_compile
    rm -f test_compile
else
    echo "❌ Test compilation failed"
    exit 1
fi

echo ""
echo "=== Step 4: Build main server ==="
if go build -o mcp-server ./cmd/server; then
    echo "✅ Main server compilation successful!"
    echo "Binary created: mcp-server"
    ls -la mcp-server
    rm -f mcp-server
else
    echo "❌ Main server compilation failed"
    exit 1
fi

echo ""
echo "🎉 All compilation tests passed!"