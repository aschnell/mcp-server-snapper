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
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "create_snapshot",
		"arguments": map[string]any{
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

	result, ok := resp["result"].(map[string]any)
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
