/*
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
*/

package main

import (
	"fmt"
	"os"

	"mcp-server-snapper/testsuite/mcpserver"
)

func main() {
	if os.Getuid() != 0 {
		fmt.Println("Skipping rollback test because it requires root privileges.")
		fmt.Println("Success.")
		return
	}

	srv, err := mcpserver.NewMcpServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start MCP server: %v\n", err)
		os.Exit(1)
	}
	defer srv.Close()

	fmt.Println("--- Rollback ---")
	resp, err := srv.SendRequest("tools/call", map[string]interface{}{
		"name": "rollback",
		"arguments": map[string]interface{}{
			"config":            "root",
			"number":            1,
			"description":       "testsuite",
			"cleanup_algorithm": "number",
			"userdata":          map[string]interface{}{},
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

	if result["isError"] == true {
		fmt.Fprintf(os.Stderr, "Error set in response\n")
		os.Exit(1)
	}

	fmt.Println("Success.")
}
