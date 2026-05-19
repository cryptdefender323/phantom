package aitools

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
	serverai "github.com/cryptdefender3232/phantom/server/ai"
	"github.com/cryptdefender3232/phantom/util/encoders"
)

// validRegistryHives lists the Windows registry root hives accepted by the implant.
var validRegistryHives = map[string]bool{
	"HKCU": true,
	"HKLM": true,
	"HKCC": true,
	"HKPD": true,
	"HKU":  true,
	"HKCR": true,
}

// validRegistryTypes maps the string names the AI uses to the protobuf uint32 constants.
var validRegistryTypes = map[string]uint32{
	"string": phantompb.RegistryTypeString,
	"binary": phantompb.RegistryTypeBinary,
	"dword":  phantompb.RegistryTypeDWORD,
	"qword":  phantompb.RegistryTypeQWORD,
}

// ---- arg structs ----

type regReadArgs struct {
	targetArgs
	Hive     string `json:"hive"`
	Path     string `json:"path"`
	Key      string `json:"key"`
	Hostname string `json:"hostname,omitempty"`
}

type regWriteArgs struct {
	targetArgs
	Hive        string  `json:"hive"`
	Path        string  `json:"path"`
	Key         string  `json:"key"`
	Type        string  `json:"type"`
	StringValue string  `json:"string_value,omitempty"`
	DWordValue  *uint32 `json:"dword_value,omitempty"`
	QWordValue  *uint64 `json:"qword_value,omitempty"`
	ByteBase64  string  `json:"byte_base64,omitempty"`
	Hostname    string  `json:"hostname,omitempty"`
}

type regListSubKeysArgs struct {
	targetArgs
	Hive     string `json:"hive"`
	Path     string `json:"path"`
	Hostname string `json:"hostname,omitempty"`
}

type regListValuesArgs struct {
	targetArgs
	Hive     string `json:"hive"`
	Path     string `json:"path"`
	Hostname string `json:"hostname,omitempty"`
}

type regCreateKeyArgs struct {
	targetArgs
	Hive     string `json:"hive"`
	Path     string `json:"path"`
	Key      string `json:"key"`
	Hostname string `json:"hostname,omitempty"`
}

type regDeleteKeyArgs struct {
	targetArgs
	Hive     string `json:"hive"`
	Path     string `json:"path"`
	Key      string `json:"key"`
	Hostname string `json:"hostname,omitempty"`
}

type regReadHiveArgs struct {
	targetArgs
	RootHive      string `json:"root_hive"`
	RequestedHive string `json:"requested_hive"`
}

// ---- result structs ----

type regReadResult struct {
	Hive  string `json:"hive"`
	Path  string `json:"path"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type regWriteResult struct {
	Hive string `json:"hive"`
	Path string `json:"path"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

type regSubKeysResult struct {
	Hive    string   `json:"hive"`
	Path    string   `json:"path"`
	Subkeys []string `json:"subkeys"`
	Count   int      `json:"count"`
}

type regValuesResult struct {
	Hive       string   `json:"hive"`
	Path       string   `json:"path"`
	ValueNames []string `json:"value_names"`
	Count      int      `json:"count"`
}

type regCreateKeyResult struct {
	Hive string `json:"hive"`
	Path string `json:"path"`
	Key  string `json:"key"`
}

type regDeleteKeyResult struct {
	Hive string `json:"hive"`
	Path string `json:"path"`
	Key  string `json:"key"`
}

type regReadHiveResult struct {
	RootHive      string `json:"root_hive"`
	RequestedHive string `json:"requested_hive"`
	ByteLen       int    `json:"byte_len"`
	SHA256        string `json:"sha256"`
	DataBase64    string `json:"data_base64"`
}

// ---- tool definitions ----

func registryToolDefinitions() []serverai.AgenticToolDefinition {
	hiveDesc := `Windows registry root hive. One of: HKCU, HKLM, HKCC, HKPD, HKU, HKCR.`
	pathDesc := `Registry key path below the hive, using backslashes, e.g. "SOFTWARE\Microsoft\Windows\CurrentVersion\Run".`
	keyDesc := `Value name within the key path.`
	hostDesc := `Optional remote hostname for registry operations on a remote machine via the implant.`

	return []serverai.AgenticToolDefinition{
		{
			Name:        "registry_read",
			Description: "Read a single Windows registry value from a session or beacon. Returns the value as a string. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":     map[string]any{"type": "string", "description": hiveDesc},
					"path":     map[string]any{"type": "string", "description": pathDesc},
					"key":      map[string]any{"type": "string", "description": keyDesc},
					"hostname": map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path", "key"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_write",
			Description: "Write a Windows registry value on a session or beacon. Supported types: string, binary, dword, qword. For binary supply byte_base64; for string supply string_value; for dword/qword supply dword_value or qword_value. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":         map[string]any{"type": "string", "description": hiveDesc},
					"path":         map[string]any{"type": "string", "description": pathDesc},
					"key":          map[string]any{"type": "string", "description": keyDesc},
					"type":         map[string]any{"type": "string", "description": `Value type: "string", "binary", "dword", or "qword".`},
					"string_value": map[string]any{"type": "string", "description": "Value to write when type is string."},
					"dword_value":  map[string]any{"type": "integer", "description": "32-bit unsigned integer to write when type is dword."},
					"qword_value":  map[string]any{"type": "integer", "description": "64-bit unsigned integer to write when type is qword."},
					"byte_base64":  map[string]any{"type": "string", "description": "Standard base64-encoded bytes to write when type is binary."},
					"hostname":     map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path", "key", "type"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_list_subkeys",
			Description: "List sub-keys under a Windows registry path on a session or beacon. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":     map[string]any{"type": "string", "description": hiveDesc},
					"path":     map[string]any{"type": "string", "description": pathDesc},
					"hostname": map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_list_values",
			Description: "List value names under a Windows registry key on a session or beacon. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":     map[string]any{"type": "string", "description": hiveDesc},
					"path":     map[string]any{"type": "string", "description": pathDesc},
					"hostname": map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_create_key",
			Description: "Create a new Windows registry key on a session or beacon. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":     map[string]any{"type": "string", "description": hiveDesc},
					"path":     map[string]any{"type": "string", "description": pathDesc},
					"key":      map[string]any{"type": "string", "description": "Name of the new sub-key to create under path."},
					"hostname": map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path", "key"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_delete_key",
			Description: "Delete a Windows registry key on a session or beacon. Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"hive":     map[string]any{"type": "string", "description": hiveDesc},
					"path":     map[string]any{"type": "string", "description": pathDesc},
					"key":      map[string]any{"type": "string", "description": "Name of the sub-key to delete under path."},
					"hostname": map[string]any{"type": "string", "description": hostDesc},
				}),
				"required":             []string{"hive", "path", "key"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "registry_read_hive",
			Description: "Dump a Windows registry hive as raw bytes from a session or beacon. Returns the hive data as base64. Useful for offline analysis (e.g. SAM, SYSTEM, SECURITY). Only works on Windows targets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": mergeSchemaProperties(map[string]any{
					"root_hive":      map[string]any{"type": "string", "description": hiveDesc + " This is the root hive that contains the requested hive."},
					"requested_hive": map[string]any{"type": "string", "description": `Sub-hive path to dump, e.g. "SAM" or "SYSTEM".`},
				}),
				"required":             []string{"root_hive", "requested_hive"},
				"additionalProperties": false,
			},
		},
	}
}

// ---- dispatch ----

func (e *executor) callRegistryTool(ctx context.Context, name string, arguments string) (string, bool, error) {
	switch strings.TrimSpace(name) {
	case "registry_read":
		var args regReadArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryRead(ctx, args)
		return result, true, err
	case "registry_write":
		var args regWriteArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryWrite(ctx, args)
		return result, true, err
	case "registry_list_subkeys":
		var args regListSubKeysArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryListSubKeys(ctx, args)
		return result, true, err
	case "registry_list_values":
		var args regListValuesArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryListValues(ctx, args)
		return result, true, err
	case "registry_create_key":
		var args regCreateKeyArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryCreateKey(ctx, args)
		return result, true, err
	case "registry_delete_key":
		var args regDeleteKeyArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryDeleteKey(ctx, args)
		return result, true, err
	case "registry_read_hive":
		var args regReadHiveArgs
		if err := decodeToolArgs(arguments, &args); err != nil {
			return "", true, err
		}
		result, err := e.callRegistryReadHive(ctx, args)
		return result, true, err
	default:
		return "", false, nil
	}
}

// ---- helpers ----

func normalizeHive(hive string) (string, error) {
	h := strings.ToUpper(strings.TrimSpace(hive))
	if !validRegistryHives[h] {
		return "", fmt.Errorf("invalid hive %q: must be one of HKCU, HKLM, HKCC, HKPD, HKU, HKCR", hive)
	}
	return h, nil
}

func normalizeRegPath(raw string) string {
	return strings.ReplaceAll(strings.TrimSpace(raw), "/", `\`)
}

// ensureWindowsTarget returns an error when the resolved target is not a Windows machine.
// It reads from the already-fetched session/beacon metadata so it costs no extra RPC.
func (e *executor) ensureWindowsTarget(ctx context.Context, sessionID, beaconID string) error {
	session, beacon, err := e.lookupTargetMetadata(ctx, sessionID, beaconID)
	if err != nil {
		return err
	}
	targetOS := ""
	if session != nil {
		targetOS = session.OS
	} else if beacon != nil {
		targetOS = beacon.OS
	}
	if !strings.EqualFold(targetOS, "windows") {
		return fmt.Errorf("registry operations require a Windows target (target OS: %q)", targetOS)
	}
	return nil
}

// ---- implementations ----

func (e *executor) callRegistryRead(ctx context.Context, args regReadArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	key := strings.TrimSpace(args.Key)
	if path == "" || key == "" {
		return "", fmt.Errorf("path and key are required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	resp, err := callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistryRead, error) {
			return e.backend.RegistryRead(callCtx, &phantompb.RegistryReadReq{
				Request:  req,
				Hive:     hive,
				Path:     path,
				Key:      key,
				Hostname: strings.TrimSpace(args.Hostname),
			})
		},
		func() *phantompb.RegistryRead { return &phantompb.RegistryRead{} },
	)
	if err != nil {
		return "", err
	}
	return marshalToolResult(regReadResult{
		Hive:  hive,
		Path:  path,
		Key:   key,
		Value: resp.Value,
	})
}

func (e *executor) callRegistryWrite(ctx context.Context, args regWriteArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	key := strings.TrimSpace(args.Key)
	if path == "" || key == "" {
		return "", fmt.Errorf("path and key are required")
	}

	typeName := strings.ToLower(strings.TrimSpace(args.Type))
	typeVal, ok := validRegistryTypes[typeName]
	if !ok {
		return "", fmt.Errorf("invalid type %q: must be one of string, binary, dword, qword", args.Type)
	}

	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	req := &phantompb.RegistryWriteReq{
		Hive:     hive,
		Path:     path,
		Key:      key,
		Type:     typeVal,
		Hostname: strings.TrimSpace(args.Hostname),
	}

	switch typeVal {
	case phantompb.RegistryTypeString:
		req.StringValue = args.StringValue
	case phantompb.RegistryTypeBinary:
		if strings.TrimSpace(args.ByteBase64) == "" {
			return "", fmt.Errorf("byte_base64 is required for binary type")
		}
		decoded, decErr := base64.StdEncoding.DecodeString(strings.TrimSpace(args.ByteBase64))
		if decErr != nil {
			decoded, decErr = base64.RawStdEncoding.DecodeString(strings.TrimSpace(args.ByteBase64))
			if decErr != nil {
				return "", fmt.Errorf("invalid byte_base64: %w", decErr)
			}
		}
		req.ByteValue = decoded
	case phantompb.RegistryTypeDWORD:
		if args.DWordValue == nil {
			return "", fmt.Errorf("dword_value is required for dword type")
		}
		req.DWordValue = *args.DWordValue
	case phantompb.RegistryTypeQWORD:
		if args.QWordValue == nil {
			return "", fmt.Errorf("qword_value is required for qword type")
		}
		req.QWordValue = *args.QWordValue
	}

	_, err = callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, commonReq *commonpb.Request) (*phantompb.RegistryWrite, error) {
			req.Request = commonReq
			return e.backend.RegistryWrite(callCtx, req)
		},
		func() *phantompb.RegistryWrite { return &phantompb.RegistryWrite{} },
	)
	if err != nil {
		return "", err
	}
	return marshalToolResult(regWriteResult{
		Hive: hive,
		Path: path,
		Key:  key,
		Type: typeName,
	})
}

func (e *executor) callRegistryListSubKeys(ctx context.Context, args regListSubKeysArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	resp, err := callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistrySubKeyList, error) {
			return e.backend.RegistryListSubKeys(callCtx, &phantompb.RegistrySubKeyListReq{
				Request:  req,
				Hive:     hive,
				Path:     path,
				Hostname: strings.TrimSpace(args.Hostname),
			})
		},
		func() *phantompb.RegistrySubKeyList { return &phantompb.RegistrySubKeyList{} },
	)
	if err != nil {
		return "", err
	}

	subkeys := resp.Subkeys
	if subkeys == nil {
		subkeys = []string{}
	}
	return marshalToolResult(regSubKeysResult{
		Hive:    hive,
		Path:    path,
		Subkeys: subkeys,
		Count:   len(subkeys),
	})
}

func (e *executor) callRegistryListValues(ctx context.Context, args regListValuesArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	resp, err := callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistryValuesList, error) {
			return e.backend.RegistryListValues(callCtx, &phantompb.RegistryListValuesReq{
				Request:  req,
				Hive:     hive,
				Path:     path,
				Hostname: strings.TrimSpace(args.Hostname),
			})
		},
		func() *phantompb.RegistryValuesList { return &phantompb.RegistryValuesList{} },
	)
	if err != nil {
		return "", err
	}

	valueNames := resp.ValueNames
	if valueNames == nil {
		valueNames = []string{}
	}
	return marshalToolResult(regValuesResult{
		Hive:       hive,
		Path:       path,
		ValueNames: valueNames,
		Count:      len(valueNames),
	})
}

func (e *executor) callRegistryCreateKey(ctx context.Context, args regCreateKeyArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	key := strings.TrimSpace(args.Key)
	if path == "" || key == "" {
		return "", fmt.Errorf("path and key are required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	_, err = callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistryCreateKey, error) {
			return e.backend.RegistryCreateKey(callCtx, &phantompb.RegistryCreateKeyReq{
				Request:  req,
				Hive:     hive,
				Path:     path,
				Key:      key,
				Hostname: strings.TrimSpace(args.Hostname),
			})
		},
		func() *phantompb.RegistryCreateKey { return &phantompb.RegistryCreateKey{} },
	)
	if err != nil {
		return "", err
	}
	return marshalToolResult(regCreateKeyResult{
		Hive: hive,
		Path: path,
		Key:  key,
	})
}

func (e *executor) callRegistryDeleteKey(ctx context.Context, args regDeleteKeyArgs) (string, error) {
	hive, err := normalizeHive(args.Hive)
	if err != nil {
		return "", err
	}
	path := normalizeRegPath(args.Path)
	key := strings.TrimSpace(args.Key)
	if path == "" || key == "" {
		return "", fmt.Errorf("path and key are required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	_, err = callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistryDeleteKey, error) {
			return e.backend.RegistryDeleteKey(callCtx, &phantompb.RegistryDeleteKeyReq{
				Request:  req,
				Hive:     hive,
				Path:     path,
				Key:      key,
				Hostname: strings.TrimSpace(args.Hostname),
			})
		},
		func() *phantompb.RegistryDeleteKey { return &phantompb.RegistryDeleteKey{} },
	)
	if err != nil {
		return "", err
	}
	return marshalToolResult(regDeleteKeyResult{
		Hive: hive,
		Path: path,
		Key:  key,
	})
}

func (e *executor) callRegistryReadHive(ctx context.Context, args regReadHiveArgs) (string, error) {
	rootHive, err := normalizeHive(args.RootHive)
	if err != nil {
		return "", err
	}
	requestedHive := normalizeRegPath(args.RequestedHive)
	if requestedHive == "" {
		return "", fmt.Errorf("requested_hive is required")
	}
	if err := e.ensureWindowsTarget(ctx, args.SessionID, args.BeaconID); err != nil {
		return "", err
	}
	target, err := e.resolveTarget(args.SessionID, args.BeaconID)
	if err != nil {
		return "", err
	}

	resp, err := callTargetRPC(
		ctx,
		target,
		func(callCtx context.Context, req *commonpb.Request) (*phantompb.RegistryReadHive, error) {
			return e.backend.RegistryReadHive(callCtx, &phantompb.RegistryReadHiveReq{
				Request:       req,
				RootHive:      rootHive,
				RequestedHive: requestedHive,
			})
		},
		func() *phantompb.RegistryReadHive { return &phantompb.RegistryReadHive{} },
	)
	if err != nil {
		return "", err
	}

	data := resp.Data
	if resp.Encoder == "gzip" {
		decoded, decErr := new(encoders.Gzip).Decode(data)
		if decErr != nil {
			return "", fmt.Errorf("failed to decompress hive data: %w", decErr)
		}
		data = decoded
	}

	return marshalToolResult(regReadHiveResult{
		RootHive:      rootHive,
		RequestedHive: requestedHive,
		ByteLen:       len(data),
		SHA256:        sha256Hex(data),
		DataBase64:    base64.StdEncoding.EncodeToString(data),
	})
}
