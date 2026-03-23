#!/usr/bin/python3

import json
from mcpserver import McpServer


mcp_server = McpServer()

print("--- Get Config ---")

response = mcp_server.send_request("tools/call", {
    "name": "get_config",
    "arguments": {
        "config": "root"
    }
})

print(json.dumps(response, indent = 2))

if not response or "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if result.get("isError"):
    raise Exception("Error set in response.")

if "structuredContent" not in result:
    raise Exception("Malformed response.")

structured_content = result["structuredContent"]

print(structured_content)

if structured_content.get("FSTYPE") != "btrfs":
    raise Exception("FSTYPE key missing or wrong value.")

print("Success.")
