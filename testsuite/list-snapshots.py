#!/usr/bin/python3
"""
Test Suite: Snapshot Management
Test Case: List Snapshots and Verify Properties

Description:
    This integration test validates the 'list_snapshots' tool of the snapper MCP server.
    It retrieves the list of snapshots for the 'root' configuration and asserts that
    the default initial snapshot (snapshot 0, description 'current') is present and
    conforms to the expected structure.

Prerequisites:
    - A standard snapper configuration named 'root' targeting '/'.
    - At least one snapshot (default snapshot 0 with description 'current').

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'list_snapshots' tool call specifying 'config': 'root'.
    3. Verify that the response contains a 'result' without 'isError'.
    4. Assert that 'structuredContent' is non-empty.
    5. Inspect the first snapshot record (index 0) and assert:
       - 'type' is 'single'
       - 'number' is 0
       - 'pre_number' is None
       - 'date' is None
       - 'description' is 'current'
       - 'cleanup_algorithm' is empty

Expected Result:
    A successful list of snapshots where the first entry accurately represents the default snapshot 0.
"""

from mcpserver import McpServer


mcp_server = McpServer()

print("--- List Snapshots ---")

response = mcp_server.send_request("tools/call", {
    "name": "list_snapshots",
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

if len(structured_content) == 0:
    raise Exception("Malformed response.")

# why again
result = structured_content["result"]

if result[0].get("type") != "single":
    raise Exception("Malformed response.")

if result[0].get("number") != 0:
    raise Exception("Malformed response.")

if result[0].get("pre_number") is not None:
    raise Exception("Malformed response.")

if result[0].get("date") is not None:
    raise Exception("Malformed response.")

if result[0].get("description") != "current":
    raise Exception("Malformed response.")

if result[0].get("cleanup_algorithm") != "":
    raise Exception("Malformed response.")

print("Success.")
