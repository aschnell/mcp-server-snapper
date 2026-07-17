/*
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

	fmt.Println("--- Get Config ---")
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "get_config",
		"arguments": map[string]any{
			"config": "root",
		},
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

	if structuredContent["FSTYPE"] != "btrfs" {
		fmt.Fprintf(os.Stderr, "FSTYPE key missing or wrong value: %v\n", structuredContent["FSTYPE"])
		os.Exit(1)
	}

	fmt.Println("Success.")
}
