#!/usr/bin/python3

import json
from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Tools ---")

response = mcp_server.send_request("tools/list")

print(json.dumps(response, indent = 2))

if not response or "result" not in response:
    raise Exception("Malformed response.")

actual_tools = { tool['name'] for tool in response["result"].get("tools", []) }

expected_tools = { "list_configs", "get_config", "set_config", "list_snapshots",
                   "create_snapshot", "delete_snapshots" }

if actual_tools != expected_tools:
    print(expected_tools)
    print(actual_tools)
    raise Exception("Wrong list of tools.")

print("Success: All expected tools are present.")
