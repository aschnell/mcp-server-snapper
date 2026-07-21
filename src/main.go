/*
 *  Copyright (c) 2026 Arvin Schnell <aschnell@suse.com>
 *  Use of this source code is governed by an MIT-style
 *  license that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version is the server version dynamically set at build time.
var Version = "0.3.0"

const snapperPath = "/usr/bin/snapper"

// Global logger file setup
var logFile *os.File

func logInfo(format string, v ...any) {
	log.Printf("INFO:root:"+format, v...)
}

func logError(format string, v ...any) {
	log.Printf("ERROR:root:"+format, v...)
}

func logDebug(format string, v ...any) {
	log.Printf("DEBUG:root:"+format, v...)
}

type Snapshot struct {
	Type             string            `json:"type"`
	Number           uint32            `json:"number"`
	PreNumber        *uint32           `json:"pre_number"`
	Date             *string           `json:"date"`
	Description      string            `json:"description"`
	CleanupAlgorithm string            `json:"cleanup_algorithm"`
	Userdata         map[string]string `json:"userdata"`
}

// Structs for MCP Tools

type ListConfigsArgs struct{}

type GetConfigArgs struct {
	Config string `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
}

type SetConfigArgs struct {
	Config string            `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
	Values map[string]string `json:"values" jsonschema:"List of key-value-pairs to set."`
}

type SetConfigOutput struct {
	Result string `json:"result"`
}

type ListSnapshotsArgs struct {
	Config string `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
}

type ListSnapshotsOutput struct {
	Result []Snapshot `json:"result" jsonschema:"Snapshots."`
}

type CreateSnapshotArgs struct {
	Config           string            `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
	Type             string            `json:"type" jsonschema:"Type for the snapshot, either 'single', 'pre' or 'post'."`
	PreNumber        uint32            `json:"pre_number" jsonschema:"Number of the corresponding pre snapshot. Required if type is 'post', otherwise ignored."`
	Description      string            `json:"description" jsonschema:"Description for the snapshot."`
	CleanupAlgorithm string            `json:"cleanup_algorithm" jsonschema:"Cleanup algorithm for the snapshot like 'number' or 'timeline'."`
	Userdata         map[string]string `json:"userdata" jsonschema:"List of key-value pairs."`
}

type CreateSnapshotOutput struct {
	Result uint32 `json:"result" jsonschema:"Number of the created snapshot."`
}

type DeleteSnapshotsArgs struct {
	Config  string   `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
	Numbers []uint32 `json:"numbers" jsonschema:"The snapshot numbers to delete."`
}

type DeleteSnapshotsOutput struct {
	Result string `json:"result"`
}

type RollbackArgs struct {
	Config           string            `json:"config" jsonschema:"Snapper config to use. Often 'root'. Use the list_configs tool to query all values."`
	Number           *uint32           `json:"number" jsonschema:"Optionally the number of the snapshot to rollback to."`
	Description      string            `json:"description" jsonschema:"Description for the new snapshot."`
	CleanupAlgorithm string            `json:"cleanup_algorithm" jsonschema:"Cleanup algorithm for the new snapshot like 'number' or 'timeline'."`
	Userdata         map[string]string `json:"userdata" jsonschema:"List of key-value pairs."`
}

type RollbackOutput struct {
	Result string `json:"result"`
}

func main() {
	log.SetOutput(os.Stderr)

	logInfo("MCP Server Snapper " + Version + " started")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run the snapper MCP server.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Initialize the MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "SnapperServer",
		Version: Version,
	}, nil)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_configs",
		Description: "Return the available snapper configs.\n" +
			":returns: Available snapper configs as a dictionary of key-value pairs with the config name as the key and the subvolume path as the value.\n" +
			":rtype: dict[str, str]",
	}, listConfigsHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "get_config",
		Description: "Return the config values of a snapper config.\n" +
			":param config: Snapper config to use. Often 'root'. Use the list_configs tool to query all values.\n" +
			":returns: Config values of a snapper config as a dictionary of key-value pairs.\n" +
			":rtype: dict[str, str]",
	}, getConfigHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "set_config",
		Description: "Set or update the configuration values for a specific snapper config.\n" +
			":param config: Snapper config to use. Often 'root'.\n" +
			":param values: List of key-value-pairs to set.",
	}, setConfigHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "list_snapshots",
		Description: "List file system snapshots using snapper.\n" +
			":param config: Snapper config to use. Often 'root'. Use the list_configs tool to query all values.\n" +
			":returns: Snapshots.\n" +
			":rtype: list[Snapshot]",
	}, listSnapshotsHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "create_snapshot",
		Description: "Create a file system snapshot using snapper.\n" +
			":param config: Snapper config to use. Often 'root'. Use the list_configs tool to query all values.\n" +
			":param type: Type for the snapshot, either 'single', 'pre' or 'post'.\n" +
			":param pre_number: Number of the corresponding pre snapshot. Required if type is 'post', otherwise ignored.\n" +
			":param description: Description for the snapshot.\n" +
			":param cleanup_algorithm: Cleanup algorithm for the snapshot like 'number' or 'timeline'.\n" +
			":param userdata: List of key-value pairs.\n" +
			":returns: Number of the created snapshot.\n" +
			":rtype: int",
	}, createSnapshotHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "delete_snapshots",
		Description: "Delete one or more file system snapshot using snapper.\n" +
			":param config: Snapper config to use. Often 'root'. Use the list_configs tool to query all values.\n" +
			":param numbers: The snapshot numbers to delete.",
	}, deleteSnapshotsHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name: "rollback",
		Description: "Rollback to a snapshot.\n" +
			":param config: Snapper config to use. Often 'root'. Use the list_configs tool to query all values.\n" +
			":param number: Optionally the number of the snapshot to rollback to.\n" +
			":param description: Description for the new snapshot.\n" +
			":param cleanup_algorithm: Cleanup algorithm for the new snapshot like 'number' or 'timeline'.\n" +
			":param userdata: List of key-value pairs.",
	}, rollbackHandler)

	// Run the server over StdioTransport
	logInfo("Starting MCP server on stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logError("Server error: %v", err)
		os.Exit(1)
	}
}

// Handlers

func listConfigsHandler(ctx context.Context, req *mcp.CallToolRequest, args ListConfigsArgs) (*mcp.CallToolResult, map[string]string, error) {
	logDebug("Received tool call: list_configs")
	resMap, err := listConfigs()
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: prettyJSON(resMap),
			},
		},
	}
	return res, resMap, nil
}

func getConfigHandler(ctx context.Context, req *mcp.CallToolRequest, args GetConfigArgs) (*mcp.CallToolResult, map[string]string, error) {
	logDebug("Received tool call: get_config with arguments: %s", compactJSON(args))
	resMap, err := getConfig(args.Config)
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: prettyJSON(resMap),
			},
		},
	}
	return res, resMap, nil
}

func setConfigHandler(ctx context.Context, req *mcp.CallToolRequest, args SetConfigArgs) (*mcp.CallToolResult, SetConfigOutput, error) {
	logDebug("Received tool call: set_config with arguments: %s", compactJSON(args))
	err := setConfig(args.Config, args.Values)
	if err != nil {
		return nil, SetConfigOutput{"success"}, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}
	return res, SetConfigOutput{Result: "success"}, nil
}

func listSnapshotsHandler(ctx context.Context, req *mcp.CallToolRequest, args ListSnapshotsArgs) (*mcp.CallToolResult, ListSnapshotsOutput, error) {
	logDebug("Received tool call: list_snapshots with arguments: %s", compactJSON(args))
	snapshots, err := listSnapshots(args.Config)
	if err != nil {
		return nil, ListSnapshotsOutput{}, err
	}
	content := make([]mcp.Content, len(snapshots))
	for i, s := range snapshots {
		content[i] = &mcp.TextContent{
			Text: prettyJSON(s),
		}
	}
	res := &mcp.CallToolResult{
		Content: content,
	}
	return res, ListSnapshotsOutput{Result: snapshots}, nil
}

func createSnapshotHandler(ctx context.Context, req *mcp.CallToolRequest, args CreateSnapshotArgs) (*mcp.CallToolResult, CreateSnapshotOutput, error) {
	logDebug("Received tool call: create_snapshot with arguments: %s", compactJSON(args))
	num, err := createSnapshot(args.Config, args.Type, args.PreNumber, args.Description, args.CleanupAlgorithm, args.Userdata)
	if err != nil {
		return nil, CreateSnapshotOutput{}, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%d", num),
			},
		},
	}
	return res, CreateSnapshotOutput{Result: num}, nil
}

func deleteSnapshotsHandler(ctx context.Context, req *mcp.CallToolRequest, args DeleteSnapshotsArgs) (*mcp.CallToolResult, DeleteSnapshotsOutput, error) {
	logDebug("Received tool call: delete_snapshots with arguments: %s", compactJSON(args))
	err := deleteSnapshots(args.Config, args.Numbers)
	if err != nil {
		return nil, DeleteSnapshotsOutput{}, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}
	return res, DeleteSnapshotsOutput{Result: "success"}, nil
}

func rollbackHandler(ctx context.Context, req *mcp.CallToolRequest, args RollbackArgs) (*mcp.CallToolResult, RollbackOutput, error) {
	logDebug("Received tool call: rollback with arguments: %s", compactJSON(args))
	err := rollback(args.Config, args.Number, args.Description, args.CleanupAlgorithm)
	if err != nil {
		return nil, RollbackOutput{}, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}
	return res, RollbackOutput{Result: "success"}, nil
}

// Core snapper business logic

func connectSystemBus() (*dbus.Conn, dbus.BusObject, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to D-Bus system bus: %w", err)
	}

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")

	return conn, obj, nil
}

func listConfigs() (map[string]string, error) {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var configs []struct {
		Name      string
		Subvolume string
		Config    map[string]string
	}
	err = obj.Call("org.opensuse.Snapper.ListConfigs", 0).Store(&configs)
	if err != nil {
		return nil, fmt.Errorf("snapper ListConfigs D-Bus call failed: %w", err)
	}

	res := make(map[string]string)
	for _, c := range configs {
		res[c.Name] = c.Subvolume
	}

	logInfo("list of snapper configs: %s", compactJSON(res))
	return res, nil
}

func getConfig(configName string) (map[string]string, error) {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var cfg struct {
		Name      string
		Subvolume string
		Config    map[string]string
	}
	err = obj.Call("org.opensuse.Snapper.GetConfig", 0, configName).Store(&cfg)
	if err != nil {
		return nil, fmt.Errorf("snapper GetConfig %q D-Bus call failed: %w", configName, err)
	}

	res := cfg.Config

	logInfo("snapper config: %s", compactJSON(res))
	return res, nil
}

func setConfig(configName string, values map[string]string) error {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = obj.Call("org.opensuse.Snapper.SetConfig", 0, configName, values).Err
	if err != nil {
		return fmt.Errorf("snapper SetConfig %q D-Bus call failed: %w", configName, err)
	}
	return nil
}

func listSnapshots(configName string) ([]Snapshot, error) {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var rawSnapshots []struct {
		Number           uint32
		Type             uint16
		PreNumber        uint32
		Timestamp        int64
		User             uint32
		Description      string
		CleanupAlgorithm string
		Userdata         map[string]string
	}
	err = obj.Call("org.opensuse.Snapper.ListSnapshots", 0, configName).Store(&rawSnapshots)
	if err != nil {
		return nil, fmt.Errorf("snapper ListSnapshots %q D-Bus call failed: %w", configName, err)
	}

	res := make([]Snapshot, 0, len(rawSnapshots))
	for _, s := range rawSnapshots {
		var tPtr *string
		if s.Timestamp != -1 {
			tStr := time.Unix(s.Timestamp, 0).Format(time.DateTime)
			tPtr = &tStr
		}

		var typeStr string
		var preNumPtr *uint32
		if s.Type == 0 {
			typeStr = "single"
		} else if s.Type == 1 {
			typeStr = "pre"
		} else if s.Type == 2 {
			typeStr = "post"
			val := s.PreNumber
			preNumPtr = &val
		}

		res = append(res, Snapshot{
			Type:             typeStr,
			Number:           s.Number,
			PreNumber:        preNumPtr,
			Date:             tPtr,
			Description:      s.Description,
			CleanupAlgorithm: s.CleanupAlgorithm,
			Userdata:         s.Userdata,
		})
	}

	logInfo("list of snapper snapshots: %s", compactJSON(res))
	return res, nil
}

func createSnapshot(configName string, typeStr string, preNumber uint32, description string, cleanupAlgorithm string, userdata map[string]string) (uint32, error) {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var number uint32

	switch typeStr {
	case "single":
		err = obj.Call("org.opensuse.Snapper.CreateSingleSnapshot", 0, configName, description, cleanupAlgorithm, userdata).Store(&number)
	case "pre":
		err = obj.Call("org.opensuse.Snapper.CreatePreSnapshot", 0, configName, description, cleanupAlgorithm, userdata).Store(&number)
	case "post":
		err = obj.Call("org.opensuse.Snapper.CreatePostSnapshot", 0, configName, preNumber, description, cleanupAlgorithm, userdata).Store(&number)
	default:
		logError("Invalid snapshot type: %s", typeStr)
		return 0, fmt.Errorf("invalid snapshot type")
	}

	if err != nil {
		return 0, fmt.Errorf("snapper create %s snapshot %q D-Bus call failed: %w", typeStr, configName, err)
	}

	logInfo("snapper number of created snapshot: %d", number)
	return number, nil
}

func deleteSnapshots(configName string, numbers []uint32) error {
	conn, obj, err := connectSystemBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = obj.Call("org.opensuse.Snapper.DeleteSnapshots", 0, configName, numbers).Err
	if err != nil {
		return fmt.Errorf("snapper DeleteSnapshots %q D-Bus call failed: %w", configName, err)
	}

	return nil
}

func rollback(configName string, number *uint32, description string, cleanupAlgorithm string) error {
	cmdArgs := []string{"--config", configName, "rollback"}
	if description != "" {
		cmdArgs = append(cmdArgs, "--description", description)
	}
	if cleanupAlgorithm != "" {
		cmdArgs = append(cmdArgs, "--cleanup-algorithm", cleanupAlgorithm)
	}
	if number != nil {
		cmdArgs = append(cmdArgs, fmt.Sprintf("%d", *number))
	}

	logInfo("%v", append([]string{snapperPath}, cmdArgs...))

	cmd := exec.Command(snapperPath, cmdArgs...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	stdout := strings.TrimSpace(stdoutBuf.String())
	stderr := strings.TrimSpace(stderrBuf.String())

	if stdout != "" {
		logInfo("STDOUT:\n%s", stdout)
	}
	if stderr != "" {
		logError("STDERR:\n%s", stderr)
	}

	if err != nil {
		var exitErr *exec.ExitError
		exitCode := 127
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		logError("Snapper error: %d", exitCode)
		if stderr != "" {
			return fmt.Errorf("snapper rollback %q command failed (exit code %d): %s: %w", configName, exitCode, stderr, err)
		}
		return fmt.Errorf("snapper rollback %q command failed (exit code %d): %w", configName, exitCode, err)
	}

	return nil
}

// Helpers

func prettyJSON(v any) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return string(bytes)
}

func compactJSON(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
