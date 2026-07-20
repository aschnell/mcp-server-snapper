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
