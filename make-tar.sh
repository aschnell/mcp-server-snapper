#!/usr/bin/bash -ex

# Get version from VERSION file
VERSION=$(cat "$(dirname "$0")/VERSION")

# Generate spec file from spec.in template
sed "s/@VERSION@/${VERSION}/g" "$(dirname "$0")/mcp-server-snapper.spec.in" > "$(dirname "$0")/package/mcp-server-snapper.spec"

# Package tarball
tar -cJvf mcp-server-snapper-${VERSION}.tar.xz --exclude='*~'	\
    --transform "s|^|mcp-server-snapper-${VERSION}/|"		\
    LICENSE README.md VERSION src/main.go go.mod go.sum vendor build.sh \
    testsuite/list-configs/main.go testsuite/get-config/main.go testsuite/tools/main.go \
    testsuite/list-snapshots/main.go testsuite/create-snapshot-1/main.go testsuite/create-snapshot-2/main.go \
    testsuite/rollback/main.go testsuite/mcpserver/mcpserver.go testsuite/README

mv mcp-server-snapper-${VERSION}.tar.xz package/
