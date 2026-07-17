package mcpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type McpServer struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	reader *bufio.Reader
}

func NewMcpServer() (*McpServer, error) {
	mcpServerBin := os.Getenv("MCPSERVER")
	if mcpServerBin == "" {
		mcpServerBin = "/usr/bin/mcp-server-snapper"
	}

	cmd := exec.Command(mcpServerBin)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	srv := &McpServer{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		reader: bufio.NewReader(stdout),
	}

	fmt.Println("--- Sending Initialize ---")
	resp, err := srv.SendRequest("initialize", map[string]any{
		"protocolVersion": "2025-11-25",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]string{
			"name":    "testsuite",
			"version": "1.0.0",
		},
	})
	if err != nil {
		srv.Close()
		return nil, fmt.Errorf("initialize failed: %w", err)
	}

	// Verify protocol version
	result, ok := resp["result"].(map[string]any)
	if !ok {
		srv.Close()
		return nil, fmt.Errorf("malformed initialize response result")
	}
	if result["protocolVersion"] != "2025-11-25" {
		srv.Close()
		return nil, fmt.Errorf("unexpected protocol version: %v", result["protocolVersion"])
	}

	return srv, nil
}

func (s *McpServer) Close() {
	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.cmd != nil {
		_ = s.cmd.Process.Kill()
		_ = s.cmd.Wait()
	}
}

func (s *McpServer) SendRequest(method string, params any) (map[string]any, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := s.stdin.Write(append(bytes, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	line, err := s.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response line: %w", err)
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty response received")
	}

	var resp map[string]any
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Print response pretty-printed
	pretty, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(pretty))

	return resp, nil
}
