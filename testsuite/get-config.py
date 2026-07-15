#!/usr/bin/python3
"""
Test Suite: Configuration Management
Test Case: Get Configuration Details

Description:
    This integration test validates the 'get_config' tool of the snapper MCP server.
    It requests configuration information for the standard 'root' configuration
    and asserts that the returned details correctly match system expectations (e.g.,
    btrfs filesystem type).

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'get_config' tool call specifying the 'root' configuration.
    3. Verify that the response contains a 'result' without 'isError'.
    4. Assert that the response includes 'structuredContent'.
    5. Verify that 'FSTYPE' key is present and has the value 'btrfs'.

Expected Result:
    Successful retrieval of configuration properties with a validated 'btrfs' filesystem type.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Get Config ---")

response = mcp_server.send_request("tools/call", {
    "name": "get_config",
    "arguments": {
        "config": "root"
    }
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

if structured_content.get("FSTYPE") != "btrfs":
    raise Exception("FSTYPE key missing or wrong value.")

print("Success.")
