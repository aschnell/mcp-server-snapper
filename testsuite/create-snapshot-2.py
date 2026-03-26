#!/usr/bin/python3

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Create Snapshot ---")

response = mcp_server.send_request("tools/call", {
    "name": "create_snapshot",
    "arguments": {
        "config": "root",
        "type": "invalid",
        "pre_number": 0,
        "description": "test",
        "cleanup": "number",
        "userdata": { "a": "b" }
    }
})

if "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if not result.get("isError"):
    raise Exception("Error not set in response.")

print("Success.")
