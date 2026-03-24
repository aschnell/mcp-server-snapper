#!/usr/bin/python3

import subprocess
import json
import sys


class McpServer:

    def __init__(self):

        self.request_id = 0

        self.process = subprocess.Popen([ "../src/mcp-server-snapper" ], stdin = subprocess.PIPE,
                                        stdout = subprocess.PIPE, stderr = sys.stderr, text = True)

        print("--- Sending Initialize ---")

        initialize_response = self.send_request("initialize", {
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
