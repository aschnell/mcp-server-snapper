#!/usr/bin/bash -ex

# Compile the main mcp-server-snapper binary
echo "Building mcp-server-snapper..."
go build -mod=vendor -o src/mcp-server-snapper src/main.go

# Compile individual test programs
echo "Building testsuite programs..."
go build -o testsuite/list-configs/list-configs testsuite/list-configs/main.go
go build -o testsuite/get-config/get-config testsuite/get-config/main.go
go build -o testsuite/tools/tools testsuite/tools/main.go
go build -o testsuite/list-snapshots/list-snapshots testsuite/list-snapshots/main.go
go build -o testsuite/create-snapshot-1/create-snapshot-1 testsuite/create-snapshot-1/main.go
go build -o testsuite/create-snapshot-2/create-snapshot-2 testsuite/create-snapshot-2/main.go
go build -o testsuite/rollback/rollback testsuite/rollback/main.go

echo "All builds completed successfully!"
