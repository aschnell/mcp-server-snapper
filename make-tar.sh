#!/usr/bin/bash -ex

# Get version from VERSION file
VERSION=$(cat "$(dirname "$0")/VERSION")

# Generate spec file from spec.in template
sed "s/@VERSION@/${VERSION}/g" "$(dirname "$0")/mcp-server-snapper.spec.in" > "$(dirname "$0")/package/mcp-server-snapper.spec"

# Package tarball
tar -cJvf mcp-server-snapper-${VERSION}.tar.xz --exclude='*~'	\
    --transform "s|^|mcp-server-snapper-${VERSION}/|"		\
    LICENSE README.md VERSION src/main.go go.mod go.sum build.sh \
    testsuite/README testsuite/*/*.go testsuite/*/README

# Package vendor tarball
tar -czf vendor.tar.gz --exclude='*~' vendor

mv mcp-server-snapper-${VERSION}.tar.xz package/
mv vendor.tar.gz package/
