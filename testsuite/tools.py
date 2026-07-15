#!/usr/bin/python3
"""
Test Suite: Discovery & Registration
Test Case: List Available Tools

Description:
    This integration test validates that the snapper MCP server correctly registers and lists
    all supported tools during discovery. It sends a 'tools/list' request and compares the
    returned list against the expected suite of tools.

Prerequisites:
    - Access to the MCP server.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'tools/list' discovery request.
    3. Verify that the response contains a 'result'.
    4. Retrieve the registered tool names and assert they exactly match the set:
       { 'list_configs', 'get_config', 'set_config', 'list_snapshots', 'create_snapshot', 'delete_snapshots', 'rollback' }

Expected Result:
    A complete list of registered tools matches the expected list exactly, with no missing or extra tools.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Tools ---")

response = mcp_server.send_request("tools/list")

if "result" not in response:
    raise Exception("Malformed response.")

expected_tools = { "list_configs", "get_config", "set_config", "list_snapshots",
                   "create_snapshot", "delete_snapshots", "rollback" }

actual_tools = { tool['name'] for tool in response["result"].get("tools", []) }

if actual_tools != expected_tools:
    print(expected_tools)
    print(actual_tools)
    raise Exception("Wrong list of tools.")

print("Success.")
