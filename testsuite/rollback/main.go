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
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "rollback",
		"arguments": map[string]any{
			"config":            "root",
			"number":            1,
			"description":       "testsuite",
			"cleanup_algorithm": "number",
			"userdata":          map[string]any{},
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

	fmt.Println("Success.")
}
