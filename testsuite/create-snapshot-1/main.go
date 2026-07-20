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
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "create_snapshot",
		"arguments": map[string]any{
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
