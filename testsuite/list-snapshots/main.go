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
