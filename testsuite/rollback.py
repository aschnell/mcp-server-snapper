#!/usr/bin/python3

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Rollback ---")

response = mcp_server.send_request("tools/call", {
    "name": "rollback",
    "arguments": {
        "config": "root",
        "number": 1,
        "description": "test",
        "cleanup": "number",
        "userdata": {}
    }
})

if "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if result.get("isError"):
    raise Exception("Error set in response.")

print("Success.")
