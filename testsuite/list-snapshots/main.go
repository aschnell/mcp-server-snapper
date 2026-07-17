/*
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

	fmt.Println("--- List Snapshots ---")
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "list_snapshots",
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

	resultsList, ok := structuredContent["result"].([]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: result is not a list\n")
		os.Exit(1)
	}

	if len(resultsList) == 0 {
		fmt.Fprintf(os.Stderr, "Empty snapshots list\n")
		os.Exit(1)
	}

	firstSnap, ok := resultsList[0].(map[string]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed snapshot object\n")
		os.Exit(1)
	}

	if firstSnap["type"] != "single" {
		fmt.Fprintf(os.Stderr, "Unexpected type: %v\n", firstSnap["type"])
		os.Exit(1)
	}

	if firstSnap["number"] != float64(0) {
		fmt.Fprintf(os.Stderr, "Unexpected number: %v\n", firstSnap["number"])
		os.Exit(1)
	}

	if firstSnap["pre_number"] != nil {
		fmt.Fprintf(os.Stderr, "Unexpected pre_number: %v\n", firstSnap["pre_number"])
		os.Exit(1)
	}

	if firstSnap["date"] != nil {
		fmt.Fprintf(os.Stderr, "Unexpected date: %v\n", firstSnap["date"])
		os.Exit(1)
	}

	if firstSnap["description"] != "current" {
		fmt.Fprintf(os.Stderr, "Unexpected description: %v\n", firstSnap["description"])
		os.Exit(1)
	}

	if firstSnap["cleanup_algorithm"] != "" {
		fmt.Fprintf(os.Stderr, "Unexpected cleanup_algorithm: %v\n", firstSnap["cleanup_algorithm"])
		os.Exit(1)
	}

	fmt.Println("Success.")
}
