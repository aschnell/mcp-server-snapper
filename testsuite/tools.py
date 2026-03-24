#!/usr/bin/python3

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Tools ---")

response = mcp_server.send_request("tools/list")

if "result" not in response:
    raise Exception("Malformed response.")

expected_tools = { "list_configs", "get_config", "set_config", "list_snapshots",
                   "create_snapshot", "delete_snapshots" }

actual_tools = { tool['name'] for tool in response["result"].get("tools", []) }

if actual_tools != expected_tools:
    print(expected_tools)
    print(actual_tools)
    raise Exception("Wrong list of tools.")

print("Success.")
