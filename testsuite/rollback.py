#!/usr/bin/python3
"""
Test Suite: System Rollback
Test Case: Snapper System Rollback

Description:
    This integration test validates the 'rollback' tool of the snapper MCP server.
    It attempts to execute a rollback to a specified snapshot number and verifies
    the tool reports a successful operation.

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.
    - A valid snapshot with number 1 (or another number specified in the arguments).
    - Appropriate administrative privileges to initiate system rollback.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'rollback' tool call with arguments:
       - config: 'root'
       - number: 1
       - description: 'testsuite'
       - cleanup_algorithm: 'number'
       - userdata: {}
    3. Verify that the response contains a 'result' without 'isError'.

Expected Result:
    A successful response indicating the rollback operation succeeded.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Rollback ---")

response = mcp_server.send_request("tools/call", {
    "name": "rollback",
    "arguments": {
        "config": "root",
        "number": 1,
        "description": "testsuite",
        "cleanup_algorithm": "number",
        "userdata": {}
    }
})

if "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if result.get("isError"):
    raise Exception("Error set in response.")

print("Success.")
