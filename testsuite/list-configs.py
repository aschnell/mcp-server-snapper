#!/usr/bin/python3

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
