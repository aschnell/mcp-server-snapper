#!/usr/bin/python3


import json
from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Tools ---")

tools_response = mcp_server.send_request("tools/list")

print(json.dumps(tools_response, indent = 2))

if not tools_response or "result" not in tools_response:
    raise Exception("Malformed response.")

actual_tools = { tool['name'] for tool in tools_response["result"].get("tools", []) }

expected_tools = { "list_configs", "get_config", "list_snapshots", "create_snapshot" }

if actual_tools != expected_tools:
    print(expected_tools)
    print(actual_tools)
    raise Exception("Wrong list of tools.")

print("Success: All expected tools are present.")
