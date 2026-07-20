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
