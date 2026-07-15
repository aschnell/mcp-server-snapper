#!/usr/bin/python3
"""
Test Utility: MCP Server Communication Harness

Description:
    This module provides the 'McpServer' helper class used across the test suite
    to spin up, communicate with, and manage the life cycle of the snapper MCP server
    subprocess. It implements the JSON-RPC 2.0 communication protocol over standard I/O (stdio).

Key Features:
    - Subprocess Management: Spawns the MCP server binary (default: '/usr/bin/mcp-server-snapper')
      or uses the location specified by the 'MCPSERVER' environment variable.
    - Handshake Protocol: Automatically executes the MCP initialization handshake protocol
      ('initialize' request with appropriate client capabilities and info) upon instantiation.
    - JSON-RPC Communication: Encapsulates sending requests, handling ID incrementation,
      and parsing incoming JSON-RPC responses.
"""

import os
import subprocess
import json
import sys


class McpServer:

    def __init__(self):

        self.request_id = 0

        mcp_server = os.environ.get('MCPSERVER', "/usr/bin/mcp-server-snapper")

        self.process = subprocess.Popen([ mcp_server ], stdin = subprocess.PIPE,
                                        stdout = subprocess.PIPE, stderr = sys.stderr,
                                        universal_newlines = True)

        print("--- Sending Initialize ---")

        self.send_request("initialize", {
            "protocolVersion": "2025-11-25",
            "capabilities": {},
            "clientInfo": { "name": "testsuite", "version": "1.0.0" }
        })


    def send_request(self, method, params = None):

        self.request_id += 1

        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }

        self.process.stdin.write(json.dumps(request) + "\n")
        self.process.stdin.flush()

        response = self.process.stdout.readline()
        if not response:
            raise Exception("Failed to read response.")

        json_response = json.loads(response)

        print(json.dumps(json_response, indent = 2))

        return json_response
