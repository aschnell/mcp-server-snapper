#!/usr/bin/bash -x

tar -cJvf mcp-server-snapper-0.1.0.tar.xz --exclude='*~'	\
    --transform 's|^|mcp-server-snapper-0.1.0/|'		\
    LICENSE README.md src/ testsuite/*.py testsuite/README

mv mcp-server-snapper-0.1.0.tar.xz package/
