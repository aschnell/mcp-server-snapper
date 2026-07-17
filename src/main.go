package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

// Snapshot represents a file system snapshot structure.
type Snapshot struct {
	Type             string            `json:"type"`
	Number           int               `json:"number"`
	PreNumber        *int              `json:"pre_number"`
	Date             *string           `json:"date"`
	Description      string            `json:"description"`
	CleanupAlgorithm string            `json:"cleanup_algorithm"`
	Userdata         map[string]string `json:"userdata"`
}

// TextContent represents a standard text block inside CallToolResult's content list.
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolResult represents the CallToolResult structure expected by MCP clients.
type ToolResult struct {
	Content           []TextContent `json:"content"`
	StructuredContent any           `json:"structuredContent,omitempty"`
	IsError           bool          `json:"isError"`
}

// RPCRequest represents an incoming JSON-RPC 2.0 request.
type RPCRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

// RPCResponse represents a successful JSON-RPC 2.0 response.
type RPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  any              `json:"result,omitempty"`
}

// RPCError represents a standard JSON-RPC 2.0 error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RPCErrorResponse represents a standard JSON-RPC 2.0 error response.
type RPCErrorResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Error   RPCError         `json:"error"`
}

// Standard JSON-RPC 2.0 error codes as specified by the JSON-RPC 2.0 specification.
const (
	ErrCodeParseError     = -32700 // Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ErrCodeInvalidRequest = -32600 // The JSON sent is not a valid Request object.
	ErrCodeMethodNotFound = -32601 // The method does not exist / is not available.
	ErrCodeInvalidParams  = -32602 // Invalid method parameter(s).
	ErrCodeInternalError  = -32603 // Internal JSON-RPC error.
)

// ToolCallParams represents the parameters for tools/call.
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// GetConfigArgs represents arguments for the get_config tool.
type GetConfigArgs struct {
	Config string `json:"config"`
}

// SetConfigArgs represents arguments for the set_config tool.
type SetConfigArgs struct {
	Config string            `json:"config"`
	Values map[string]string `json:"values"`
}

// ListSnapshotsArgs represents arguments for the list_snapshots tool.
type ListSnapshotsArgs struct {
	Config string `json:"config"`
}

// CreateSnapshotArgs represents arguments for the create_snapshot tool.
type CreateSnapshotArgs struct {
	Config           string            `json:"config"`
	Type             string            `json:"type"`
	PreNumber        int               `json:"pre_number"`
	Description      string            `json:"description"`
	CleanupAlgorithm string            `json:"cleanup_algorithm"`
	Userdata         map[string]string `json:"userdata"`
}

// DeleteSnapshotsArgs represents arguments for the delete_snapshots tool.
type DeleteSnapshotsArgs struct {
	Config  string `json:"config"`
	Numbers []int  `json:"numbers"`
}

// RollbackArgs represents arguments for the rollback tool.
type RollbackArgs struct {
	Config           string            `json:"config"`
	Number           *int              `json:"number"`
	Description      string            `json:"description"`
	CleanupAlgorithm string            `json:"cleanup_algorithm"`
	Userdata         map[string]string `json:"userdata"`
}

// Version is the server version dynamically set at build time.
var Version = "0.3.0"

// Global logger file setup
var logFile *os.File

func initLogger() {
	var err error
	logFile, err = os.OpenFile("mcp-server-snapper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(logFile)
	}
}

func logInfo(format string, v ...any) {
	log.Printf("INFO:root:"+format, v...)
}

func logError(format string, v ...any) {
	log.Printf("ERROR:root:"+format, v...)
}

func logDebug(format string, v ...any) {
	log.Printf("DEBUG:root:"+format, v...)
}

func main() {
	initLogger()
	logInfo("Server started")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run the snapper MCP server.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logError("Error reading stdin: %v", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		logDebug("Received message raw line: %s", line)

		var req RPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			logError("Failed to parse request JSON: %v", err)
			resp := makeErrorResponse(nil, ErrCodeParseError, "Parse error")
			fmt.Println(resp)
			continue
		}

		// Process the request
		respLine := handleRPCRequest(&req)
		if respLine != "" {
			fmt.Println(respLine)
		}
	}
}

func handleRPCRequest(req *RPCRequest) string {
	switch req.Method {
	case "initialize":
		return handleInitialize(req)
	case "tools/list":
		return handleToolsList(req)
	case "tools/call":
		return handleToolsCall(req)
	default:
		// If it's a notification, do not send response
		if req.ID == nil {
			logDebug("Ignored notification: %s", req.Method)
			return ""
		}
		logError("Method not found: %s", req.Method)
		return makeErrorResponse(req.ID, ErrCodeMethodNotFound, "Method not found")
	}
}

func handleInitialize(req *RPCRequest) string {
	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]any{
			"protocolVersion": "2025-11-25",
			"capabilities": map[string]any{
				"experimental": map[string]any{},
				"prompts": map[string]any{
					"listChanged": false,
				},
				"resources": map[string]any{
					"subscribe":   false,
					"listChanged": false,
				},
				"tools": map[string]any{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]string{
				"name":    "SnapperServer",
				"version": Version,
			},
		},
	}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}

// Static definition of tools matching python tools/list response schemas exactly.
const toolsListJSON = `{
  "tools": [
    {
      "name": "list_configs",
      "description": "\nReturn the available snapper configs.\n:returns: Available snapper configs as a dictionary of key-value pairs with the config\n          name as the key and the subvolume path as the value.\n:rtype: dict[str, str]\n",
      "inputSchema": {
        "properties": {},
        "title": "list_configsArguments",
        "type": "object"
      },
      "outputSchema": {
        "additionalProperties": {
          "type": "string"
        },
        "title": "list_configsDictOutput",
        "type": "object"
      }
    },
    {
      "name": "get_config",
      "description": "\nReturn the config values of a snapper config.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:returns: Config values of a snapper config as a dictionary of key-value pairs.\n:rtype: dict[str, str]\n",
      "inputSchema": {
        "properties": {
          "config": {
            "title": "Config",
            "type": "string"
          }
        },
        "required": [
          "config"
        ],
        "title": "get_configArguments",
        "type": "object"
      },
      "outputSchema": {
        "additionalProperties": {
          "type": "string"
        },
        "title": "get_configDictOutput",
        "type": "object"
      }
    },
    {
      "name": "set_config",
      "description": "\nList the configuration values of a snapper config.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:param values: List of key-value-pairs to set.\n",
      "inputSchema": {
        "properties": {
          "config": {
            "title": "Config",
            "type": "string"
          },
          "values": {
            "additionalProperties": {
              "type": "string"
            },
            "title": "Values",
            "type": "object"
          }
        },
        "required": [
          "config",
          "values"
        ],
        "title": "set_configArguments",
        "type": "object"
      },
      "outputSchema": {
        "properties": {
          "result": {
            "title": "Result",
            "type": "null"
          }
        },
        "required": [
          "result"
        ],
        "title": "set_configOutput",
        "type": "object"
      }
    },
    {
      "name": "list_snapshots",
      "description": "\nList file system snapshots using snapper.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:returns: Snapshots.\n:rtype: list[Snapshot]\n",
      "inputSchema": {
        "properties": {
          "config": {
            "title": "Config",
            "type": "string"
          }
        },
        "required": [
          "config"
        ],
        "title": "list_snapshotsArguments",
        "type": "object"
      },
      "outputSchema": {
        "$defs": {
          "Snapshot": {
            "properties": {
              "cleanup_algorithm": {
                "title": "Cleanup Algorithm",
                "type": "string"
              },
              "date": {
                "anyOf": [
                  {
                    "type": "string"
                  },
                  {
                    "type": "null"
                  }
                ],
                "default": null,
                "title": "Date"
              },
              "description": {
                "title": "Description",
                "type": "string"
              },
              "number": {
                "title": "Number",
                "type": "integer"
              },
              "pre_number": {
                "anyOf": [
                  {
                    "type": "integer"
                  },
                  {
                    "type": "null"
                  }
                ],
                "default": null,
                "title": "Pre Number"
              },
              "type": {
                "title": "Type",
                "type": "string"
              },
              "userdata": {
                "additionalProperties": {
                  "type": "string"
                },
                "title": "Userdata",
                "type": "object"
              }
            },
            "required": [
              "type",
              "number",
              "description",
              "cleanup_algorithm",
              "userdata"
            ],
            "title": "Snapshot",
            "type": "object"
          }
        },
        "properties": {
          "result": {
            "items": {
              "$ref": "#/$defs/Snapshot"
            },
            "title": "Result",
            "type": "array"
          }
        },
        "required": [
          "result"
        ],
        "title": "list_snapshotsOutput",
        "type": "object"
      }
    },
    {
      "name": "create_snapshot",
      "description": "\nCreate a file system snapshot using snapper.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:param type: Type for the snapshot, either 'single', 'pre' or 'post'.\n:param pre_number: Number of the corresponding pre snapshot. Required if type is 'post',\n       otherwise ignored.\n:param description: Description for the snapshot.\n:param cleanup_algorithm: Cleanup algorithm for the snapshot like 'number' or 'timeline'.\n:param userdata: List of key-value pairs.\n:returns: Number of the created snapshot.\n:rtype: int\n",
      "inputSchema": {
        "properties": {
          "cleanup_algorithm": {
            "title": "Cleanup Algorithm",
            "type": "string"
          },
          "config": {
            "title": "Config",
            "type": "string"
          },
          "description": {
            "title": "Description",
            "type": "string"
          },
          "pre_number": {
            "title": "Pre Number",
            "type": "integer"
          },
          "type": {
            "title": "Type",
            "type": "string"
          },
          "userdata": {
            "additionalProperties": {
              "type": "string"
            },
            "title": "Userdata",
            "type": "object"
          }
        },
        "required": [
          "config",
          "type",
          "pre_number",
          "description",
          "cleanup_algorithm",
          "userdata"
        ],
        "title": "create_snapshotArguments",
        "type": "object"
      },
      "outputSchema": {
        "properties": {
          "result": {
            "title": "Result",
            "type": "integer"
          }
        },
        "required": [
          "result"
        ],
        "title": "create_snapshotOutput",
        "type": "object"
      }
    },
    {
      "name": "delete_snapshots",
      "description": "\nDelete one or more file system snapshot using snapper.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:param numbers: The snapshot numbers to delete.\n",
      "inputSchema": {
        "properties": {
          "config": {
            "title": "Config",
            "type": "string"
          },
          "numbers": {
            "items": {
              "type": "integer"
            },
            "title": "Numbers",
            "type": "array"
          }
        },
        "required": [
          "config",
          "numbers"
        ],
        "title": "delete_snapshotsArguments",
        "type": "object"
      },
      "outputSchema": {
        "properties": {
          "result": {
            "title": "Result",
            "type": "null"
          }
        },
        "required": [
          "result"
        ],
        "title": "delete_snapshotsOutput",
        "type": "object"
      }
    },
    {
      "name": "rollback",
      "description": "\nRollback to a snapshot.\n:param config: Snapper config to use. Often 'root'. Use the list_configs tool to\n       query all values.\n:param number: Optionally the number of the snapshot to rollback to.\n:param description: Description for the new snapshot.\n:param cleanup_algorithm: Cleanup algorithm for the new snapshot like 'number' or 'timeline'.\n:param userdata: List of key-value pairs.\n",
      "inputSchema": {
        "properties": {
          "cleanup_algorithm": {
            "title": "Cleanup Algorithm",
            "type": "string"
          },
          "config": {
            "title": "Config",
            "type": "string"
          },
          "description": {
            "title": "Description",
            "type": "string"
          },
          "number": {
            "anyOf": [
              {
                "type": "integer"
              },
              {
                "type": "null"
              }
            ],
            "title": "Number"
          },
          "userdata": {
            "additionalProperties": {
              "type": "string"
            },
            "title": "Userdata",
            "type": "object"
          }
        },
        "required": [
          "config",
          "number",
          "description",
          "cleanup_algorithm",
          "userdata"
        ],
        "title": "rollbackArguments",
        "type": "object"
      },
      "outputSchema": {
        "properties": {
          "result": {
            "title": "Result",
            "type": "null"
          }
        },
        "required": [
          "result"
        ],
        "title": "rollbackOutput",
        "type": "object"
      }
    }
  ]
}`

func handleToolsList(req *RPCRequest) string {
	var toolsObj map[string]any
	_ = json.Unmarshal([]byte(toolsListJSON), &toolsObj)

	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  toolsObj,
	}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}

func handleToolsCall(req *RPCRequest) string {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		logError("Failed to unmarshal tool call params: %v", err)
		return makeErrorResponse(req.ID, ErrCodeInvalidParams, "Invalid params")
	}

	var val any
	var err error

	switch params.Name {
	case "list_configs":
		val, err = listConfigs()
	case "get_config":
		var args GetConfigArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = getConfig(args.Config)
	case "set_config":
		var args SetConfigArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = setConfig(args.Config, args.Values)
	case "list_snapshots":
		var args ListSnapshotsArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = listSnapshots(args.Config)
	case "create_snapshot":
		var args CreateSnapshotArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = createSnapshot(args.Config, args.Type, args.PreNumber, args.Description, args.CleanupAlgorithm, args.Userdata)
	case "delete_snapshots":
		var args DeleteSnapshotsArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = deleteSnapshots(args.Config, args.Numbers)
	case "rollback":
		var args RollbackArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return makeToolErrorResponse(req.ID, params.Name, fmt.Errorf("invalid arguments: %w", err))
		}
		val, err = rollback(args.Config, args.Number, args.Description, args.CleanupAlgorithm)
	default:
		return makeErrorResponse(req.ID, ErrCodeMethodNotFound, "Tool not found")
	}

	if err != nil {
		logError("Error executing tool %s: %v", params.Name, err)
		return makeToolErrorResponse(req.ID, params.Name, err)
	}

	return makeSuccessResponse(req.ID, params.Name, val)
}

func listConfigs() (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
	var configs []struct {
		Name      string
		Subvolume string
		Config    map[string]string
	}
	err = obj.Call("org.opensuse.Snapper.ListConfigs", 0).Store(&configs)
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}

	res := make(map[string]string)
	for _, c := range configs {
		res[c.Name] = c.Subvolume
	}
	logInfo("list of snapper configs: %s", compactJSON(res))
	return res, nil
}

func getConfig(configName string) (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
	var cfg struct {
		Name      string
		Subvolume string
		Config    map[string]string
	}
	err = obj.Call("org.opensuse.Snapper.GetConfig", 0, configName).Store(&cfg)
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}

	res := make(map[string]string)
	for k, v := range cfg.Config {
		res[k] = v
	}
	logInfo("snapper config: %s", compactJSON(res))
	return res, nil
}

func setConfig(configName string, values map[string]string) (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
	err = obj.Call("org.opensuse.Snapper.SetConfig", 0, configName, values).Err
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	return nil, nil
}

func listSnapshots(configName string) (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
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
		return nil, fmt.Errorf("snapper error")
	}

	res := make([]Snapshot, 0, len(rawSnapshots))
	for _, s := range rawSnapshots {
		var tPtr *string
		if s.Timestamp != -1 {
			tStr := time.Unix(s.Timestamp, 0).Format("2006-01-02 15:04:05")
			tPtr = &tStr
		}

		var typeStr string
		var preNumPtr *int
		if s.Type == 0 {
			typeStr = "single"
		} else if s.Type == 1 {
			typeStr = "pre"
		} else if s.Type == 2 {
			typeStr = "post"
			val := int(s.PreNumber)
			preNumPtr = &val
		}

		userdata := s.Userdata
		if userdata == nil {
			userdata = make(map[string]string)
		}

		res = append(res, Snapshot{
			Type:             typeStr,
			Number:           int(s.Number),
			PreNumber:        preNumPtr,
			Date:             tPtr,
			Description:      s.Description,
			CleanupAlgorithm: s.CleanupAlgorithm,
			Userdata:         userdata,
		})
	}

	logInfo("list of snapper snapshots: %s", compactJSON(res))
	return res, nil
}

func createSnapshot(configName string, typeStr string, preNumber int, description string, cleanupAlgorithm string, userdata map[string]string) (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
	var number uint32

	if userdata == nil {
		userdata = make(map[string]string)
	}

	switch typeStr {
	case "single":
		err = obj.Call("org.opensuse.Snapper.CreateSingleSnapshot", 0, configName, description, cleanupAlgorithm, userdata).Store(&number)
	case "pre":
		err = obj.Call("org.opensuse.Snapper.CreatePreSnapshot", 0, configName, description, cleanupAlgorithm, userdata).Store(&number)
	case "post":
		err = obj.Call("org.opensuse.Snapper.CreatePostSnapshot", 0, configName, uint32(preNumber), description, cleanupAlgorithm, userdata).Store(&number)
	default:
		logError("Invalid snapshot type: %s", typeStr)
		return nil, fmt.Errorf("invalid snapshot type")
	}

	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}

	logInfo("snapper number of created snapshot: %d", number)
	return int(number), nil
}

func deleteSnapshots(configName string, numbers []int) (any, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	defer conn.Close()

	obj := conn.Object("org.opensuse.Snapper", "/org/opensuse/Snapper")
	dbusNums := make([]uint32, len(numbers))
	for i, n := range numbers {
		dbusNums[i] = uint32(n)
	}
	err = obj.Call("org.opensuse.Snapper.DeleteSnapshots", 0, configName, dbusNums).Err
	if err != nil {
		return nil, fmt.Errorf("snapper error")
	}
	return nil, nil
}

const snapperPath = "/usr/bin/snapper"

func rollback(configName string, number *int, description string, cleanupAlgorithm string) (any, error) {
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
		return nil, fmt.Errorf("snapper error")
	}

	return nil, nil
}

func makeSuccessResponse(id *json.RawMessage, name string, val any) string {
	var toolRes ToolResult
	toolRes.IsError = false

	if val == nil {
		toolRes.Content = []TextContent{}
		toolRes.StructuredContent = map[string]any{"result": nil}
	} else {
		switch v := val.(type) {
		case map[string]string:
			toolRes.StructuredContent = v
			toolRes.Content = []TextContent{
				{
					Type: "text",
					Text: prettyJSON(v),
				},
			}
		case []Snapshot:
			toolRes.StructuredContent = map[string]any{"result": v}
			toolRes.Content = make([]TextContent, len(v))
			for i, s := range v {
				toolRes.Content[i] = TextContent{
					Type: "text",
					Text: prettyJSON(s),
				}
			}
		case int:
			toolRes.StructuredContent = map[string]any{"result": v}
			toolRes.Content = []TextContent{
				{
					Type: "text",
					Text: fmt.Sprintf("%d", v),
				},
			}
		default:
			toolRes.StructuredContent = v
			toolRes.Content = []TextContent{
				{
					Type: "text",
					Text: prettyJSON(v),
				},
			}
		}
	}

	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  toolRes,
	}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}

func makeToolErrorResponse(id *json.RawMessage, name string, err error) string {
	res := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: ToolResult{
			Content: []TextContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error executing tool %s: %s", name, err.Error()),
				},
			},
			IsError: true,
		},
	}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}

func makeErrorResponse(id *json.RawMessage, code int, message string) string {
	res := RPCErrorResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: RPCError{
			Code:    code,
			Message: message,
		},
	}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}

func prettyJSON(v any) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return string(bytes)
}

func compactJSON(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
