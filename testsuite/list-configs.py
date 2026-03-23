#!/usr/bin/python3

import json
from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Configs ---")

response = mcp_server.send_request("tools/call", {
    "name": "list_configs"
})

print(json.dumps(response, indent = 2))

if not response or "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if result.get("isError"):
    raise Exception("Error set in response.")

if "content" not in result:
    raise Exception("Malformed response.")

content = result["content"]

if len(content) != 1:
    raise Exception("Malformed response.")

if content[0].get('type') != "text":
    raise Exception("Wrong type in content.")

expected_text = "List of snapper configs as table with header: ['Config | Subvolume', 'root | /']"

actual_text = content[0].get('text')

if actual_text != expected_text:
    print(expected_text)
    print(actual_text)
    raise Exception("Wrong text.")

print("Success.")
