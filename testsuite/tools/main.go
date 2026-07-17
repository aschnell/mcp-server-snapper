/*
Test Suite: Discovery & Registration
Test Case: List Available Tools

Description:
    This integration test validates that the snapper MCP server correctly registers and lists
    all supported tools during discovery. It sends a 'tools/list' request and compares the
    returned list against the expected suite of tools.

Prerequisites:
    - Access to the MCP server.

Test Steps:
    1. Initialize the MCP server connection.
    2. Dispatch a 'tools/list' discovery request.
    3. Verify that the response contains a 'result'.
    4. Retrieve the registered tool names and assert they exactly match the set:
       { 'list_configs', 'get_config', 'set_config', 'list_snapshots', 'create_snapshot', 'delete_snapshots', 'rollback' }

Expected Result:
    A complete list of registered tools matches the expected list exactly, with no missing or extra tools.
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

	fmt.Println("--- Listing Tools ---")
	resp, err := srv.SendRequest("tools/list", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
		os.Exit(1)
	}

	result, ok := resp["result"].(map[string]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: result missing\n")
		os.Exit(1)
	}

	toolsList, ok := result["tools"].([]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed response: tools missing\n")
		os.Exit(1)
	}

	expectedTools := map[string]bool{
		"list_configs":     true,
		"get_config":       true,
		"set_config":       true,
		"list_snapshots":   true,
		"create_snapshot":  true,
		"delete_snapshots": true,
		"rollback":         true,
	}

	actualTools := make(map[string]bool)
	for _, tVal := range toolsList {
		tObj, ok := tVal.(map[string]any)
		if !ok {
			fmt.Fprintf(os.Stderr, "Malformed tool object\n")
			os.Exit(1)
		}
		name, ok := tObj["name"].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "Tool name is not string\n")
			os.Exit(1)
		}
		actualTools[name] = true
	}

	for expected := range expectedTools {
		if !actualTools[expected] {
			fmt.Fprintf(os.Stderr, "Missing expected tool: %s\n", expected)
			os.Exit(1)
		}
	}

	if len(actualTools) != len(expectedTools) {
		fmt.Fprintf(os.Stderr, "Wrong list of tools. Expected: %v, got: %v\n", expectedTools, actualTools)
		os.Exit(1)
	}

	fmt.Println("Success.")
}
