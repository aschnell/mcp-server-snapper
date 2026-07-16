#!/usr/bin/bash -x

tar -cJvf mcp-server-snapper-0.3.0.tar.xz --exclude='*~'	\
    --transform 's|^|mcp-server-snapper-0.3.0/|'		\
    LICENSE README.md src/main.go go.mod go.sum vendor build.sh \
    testsuite/list-configs/main.go testsuite/get-config/main.go testsuite/tools/main.go \
    testsuite/list-snapshots/main.go testsuite/create-snapshot-1/main.go testsuite/create-snapshot-2/main.go \
    testsuite/rollback/main.go testsuite/mcpserver/mcpserver.go testsuite/README

mv mcp-server-snapper-0.3.0.tar.xz package/
