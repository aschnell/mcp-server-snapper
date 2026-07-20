package main

import (
	"fmt"
	"os"
	"strings"

	"mcp-server-snapper/testsuite/mcpserver"
)

func main() {
	srv, err := mcpserver.NewMcpServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start MCP server: %v\n", err)
		os.Exit(1)
	}
	defer srv.Close()

	fmt.Println("--- Create Snapshot Non-Existent Config (Detailed Error) ---")
	resp, err := srv.SendRequest("tools/call", map[string]any{
		"name": "create_snapshot",
		"arguments": map[string]any{
			"config":            "non_existent_config_xyz",
			"type":              "single",
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
		fmt.Fprintf(os.Stderr, "Error not set in response for non-existent config\n")
		os.Exit(1)
	}

	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		fmt.Fprintf(os.Stderr, "Content missing or empty in error response\n")
		os.Exit(1)
	}

	contentMap, ok := content[0].(map[string]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Malformed content item\n")
		os.Exit(1)
	}

	errMsg, ok := contentMap["text"].(string)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error message text missing\n")
		os.Exit(1)
	}

	expectedPrefix := `snapper create single snapshot "non_existent_config_xyz" D-Bus call failed`
	if !strings.Contains(errMsg, expectedPrefix) {
		fmt.Fprintf(os.Stderr, "Expected error message to contain %q, but got: %q\n", expectedPrefix, errMsg)
		os.Exit(1)
	}

	fmt.Printf("Success.")
}
