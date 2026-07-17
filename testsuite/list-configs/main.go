/*
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
*/

package main

import (
	"fmt"
	"os"

	"mcp-server-snapper/testsuite/mcpserver"
)

func main() {
	srv, err := mcpserver.NewMcpServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start MCP server: %v\n", err)
		os.Exit(1)
	}
	defer srv.Close()

	fmt.Println("--- Listing Configs ---")
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "list_configs",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
		os.Exit(1)
	}

	result, ok := resp["result"].(map[string]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: result missing\n")
		os.Exit(1)
	}

	if result["isError"] == true {
		fmt.Fprintf(os.Stderr, "Error set in response\n")
		os.Exit(1)
	}

	structuredContent, ok := result["structuredContent"].(map[string]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: structuredContent missing\n")
		os.Exit(1)
	}

	if structuredContent["root"] != "/" {
		fmt.Fprintf(os.Stderr, "root key missing or wrong value: %v\n", structuredContent["root"])
		os.Exit(1)
	}

	fmt.Println("Success.")
}
