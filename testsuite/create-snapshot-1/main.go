/*
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

	fmt.Println("--- Create Snapshot ---")
	resp, err := srv.SendRequest("tools/call", map[string]interface{}{
		"name": "create_snapshot",
		"arguments": map[string]interface{}{
			"config":            "root",
			"type":              "single",
			"pre_number":        0,
			"description":       "testsuite",
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

	if result["isError"] == true {
		fmt.Fprintf(os.Stderr, "Error set in response\n")
		os.Exit(1)
	}

	structuredContent, ok := result["structuredContent"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: structuredContent missing\n")
		os.Exit(1)
	}

	numVal, ok := structuredContent["result"]
	if !ok {
		fmt.Fprintf(os.Stderr, "StructuredContent result missing\n")
		os.Exit(1)
	}

	floatNum, ok := numVal.(float64)
	if !ok {
		fmt.Fprintf(os.Stderr, "Wrong type of result: expected number, got %T\n", numVal)
		os.Exit(1)
	}

	snapshotNum := int(floatNum)
	fmt.Printf("Created snapshot: %d\n", snapshotNum)
	fmt.Println("Success.")
}
