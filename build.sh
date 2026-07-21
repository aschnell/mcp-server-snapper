#!/usr/bin/bash -ex

# Get version from VERSION file
VERSION=$(cat "$(dirname "$0")/VERSION")

# Compile the main mcp-server-snapper binary
echo "Building mcp-server-snapper with version ${VERSION}..."
go build -mod=vendor -buildmode=pie -ldflags "-X main.Version=${VERSION}" -o src/mcp-server-snapper src/main.go

# Compile individual test programs
echo "Building testsuite programs..."
for prog in create-snapshot-1 create-snapshot-2 create-snapshot-3 get-config list-configs list-snapshots rollback tools; do
    go build -mod=vendor -buildmode=pie -o testsuite/$prog/$prog testsuite/$prog/main.go
done

echo "All builds completed successfully!"
