cat > scripts/build-readmegen.sh << 'EOF'
#!/bin/bash

# Build the readme generator
echo "Building readmegen..."

cd "$(dirname "$0")/.." || exit 1

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build
go build -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE" \
    -o bin/readmegen cmd/readmegen/main.go

if [ $? -eq 0 ]; then
    echo "✅ Built bin/readmegen (version: $VERSION)"
else
    echo "❌ Build failed"
    exit 1
fi
EOF

chmod +x scripts/build-readmegen.sh