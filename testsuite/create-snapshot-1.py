#!/usr/bin/python3
"""
Test Suite: Snapshot Management
Test Case: Successful Single Snapshot Creation

Description:
    This integration test verifies the successful creation of a new snapper snapshot
    via the 'create_snapshot' tool of the MCP server. It connects to the MCP server,
    sends a 'tools/call' request specifying a 'single' type snapshot on the default
    'root' configuration, and validates that a valid snapshot number is returned.

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.
    - Appropriate permissions to execute snapper operations.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'create_snapshot' tool call with arguments:
       - config: 'root'
       - type: 'single'
       - pre_number: 0
       - description: 'testsuite'
       - cleanup_algorithm: 'number'
       - userdata: {"a": "b"}
    3. Assert the server's response contains a valid JSON-RPC 'result'.
    4. Assert that 'isError' is either absent or false in the response.
    5. Retrieve the 'structuredContent' and assert that the returned snapshot identifier
       is a valid integer.

Expected Result:
    The response contains no errors and yields a positive integer representing the ID
    of the newly created snapshot.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- Create Snapshot ---")

response = mcp_server.send_request("tools/call", {
    "name": "create_snapshot",
    "arguments": {
        "config": "root",
        "type": "single",
        "pre_number": 0,
        "description": "testsuite",
        "cleanup_algorithm": "number",
        "userdata": { "a": "b" }
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

number = structured_content["result"]

print(number)

if not isinstance(number, int):
    raise Exception("Wrong type of result.")

print("Success.")
