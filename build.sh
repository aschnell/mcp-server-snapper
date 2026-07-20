#!/usr/bin/bash -ex

# Get version from VERSION file
VERSION=$(cat "$(dirname "$0")/VERSION")

# Compile the main mcp-server-snapper binary
echo "Building mcp-server-snapper with version ${VERSION}..."
go build -mod=vendor -buildmode=pie -ldflags "-X main.Version=${VERSION}" -o src/mcp-server-snapper src/main.go

# Compile individual test programs
echo "Building testsuite programs..."
go build -mod=vendor -buildmode=pie -o testsuite/list-configs/list-configs testsuite/list-configs/main.go
go build -mod=vendor -buildmode=pie -o testsuite/get-config/get-config testsuite/get-config/main.go
go build -mod=vendor -buildmode=pie -o testsuite/tools/tools testsuite/tools/main.go
go build -mod=vendor -buildmode=pie -o testsuite/list-snapshots/list-snapshots testsuite/list-snapshots/main.go
go build -mod=vendor -buildmode=pie -o testsuite/create-snapshot-1/create-snapshot-1 testsuite/create-snapshot-1/main.go
go build -mod=vendor -buildmode=pie -o testsuite/create-snapshot-2/create-snapshot-2 testsuite/create-snapshot-2/main.go
go build -mod=vendor -buildmode=pie -o testsuite/rollback/rollback testsuite/rollback/main.go

echo "All builds completed successfully!"
