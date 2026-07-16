/*
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

	fmt.Println("--- Create Snapshot Invalid ---")
	resp, err := srv.SendRequest("tools/call", map[string]interface{}{
		"name": "create_snapshot",
		"arguments": map[string]interface{}{
			"config":            "root",
			"type":              "invalid",
			"pre_number":        0,
			"description":       "test",
			"cleanup_algorithm": "number",
			"userdata":          map[string]string{"a": "b"},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
		os.Exit(1)
	}

	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: result missing\n")
		os.Exit(1)
	}

	if result["isError"] != true {
		fmt.Fprintf(os.Stderr, "Error not set in response\n")
		os.Exit(1)
	}

	fmt.Println("Success.")
}
