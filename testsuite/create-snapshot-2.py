#!/usr/bin/python3
"""
Test Suite: Snapshot Management
Test Case: Error Handling on Invalid Snapshot Type

Description:
    This integration test verifies that the MCP server robustly handles invalid input
    during snapshot creation. It sends a 'tools/call' request to the 'create_snapshot'
    tool with an unsupported snapshot type ('invalid') and verifies that the server
    correctly propagates an error response instead of failing silently.

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'create_snapshot' tool call with arguments containing:
       - config: 'root'
       - type: 'invalid' (unsupported value)
       - other standard parameters.
    3. Assert the server's response contains a valid JSON-RPC 'result'.
    4. Validate that 'isError' is present and evaluated to True in the result, ensuring
       proper server-side input validation and error reporting.

Expected Result:
    The response succeeds at the RPC level but includes an 'isError' flag set to True,
    confirming the server rejected the invalid arguments gracefully.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Create Snapshot ---")

response = mcp_server.send_request("tools/call", {
    "name": "create_snapshot",
    "arguments": {
        "config": "root",
        "type": "invalid",
        "pre_number": 0,
        "description": "test",
        "cleanup_algorithm": "number",
        "userdata": { "a": "b" }
    }
})

if "result" not in response:
    raise Exception("Malformed response.")

result = response["result"]

if not result.get("isError"):
    raise Exception("Error not set in response.")

print("Success.")
