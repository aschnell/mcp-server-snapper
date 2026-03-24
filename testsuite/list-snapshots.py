#!/usr/bin/python3

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

if result[0].get("cleanup") != "":
    raise Exception("Malformed response.")

print("Success.")
