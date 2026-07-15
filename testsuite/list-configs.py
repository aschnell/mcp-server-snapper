#!/usr/bin/python3
"""
Test Suite: Configuration Management
Test Case: List Snapper Configurations

Description:
    This integration test verifies the 'list_configs' tool of the snapper MCP server.
    It calls the tool to list all snapper configurations defined on the system and
    verifies that the default 'root' configuration is present and mapped to '/'.

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'list_configs' tool call.
    3. Verify that the response contains a 'result' without 'isError'.
    4. Retrieve 'structuredContent' and assert the 'root' key exists and is mapped to '/'.

Expected Result:
    A successful list of configurations, containing at least the standard 'root' configuration.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Listing Configs ---")

response = mcp_server.send_request("tools/call", {
    "name": "list_configs"
})

if "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if result.get("isError"):
    raise Exception("Error set in response.")

if "structuredContent" not in result:
    raise Exception("Malformed response.")

structured_content = result["structuredContent"]

print(structured_content)

if structured_content.get('root') != "/":
    raise Exception("root key missing or wrong value.")

print("Success.")
