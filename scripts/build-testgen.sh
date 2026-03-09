cat > scripts/build-testgen.sh << 'EOF'
#!/bin/bash

# Build the test generator
echo "Building testgen..."

cd "$(dirname "$0")/.." || exit 1

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build
go build -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE" \
    -o bin/testgen.exe cmd/testgen/main.go

if [ $? -eq 0 ]; then
    echo "✅ Built bin/testgen.exe (version: $VERSION)"
else
    echo "❌ Build failed"
    exit 1
fi
EOF

chmod +x scripts/build-testgen.sh