#!/usr/bin/python3

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
        "cleanup": "number",
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
