// Package internal contains internal utilities and logic for the script manager CLI tool.
package internal

// ScriptInfo represents metadata about a script.
type ScriptInfo struct {
	ScriptName string `json:"script_name"` // The name of the script (without extension).
	ScriptExt  string `json:"script_ext"`  // The file extension of the script.
	ServerPath string `json:"server_path"` // The path to the script on the server.
	ServerAddr string `json:"server_addr"` // The address of the server hosting the script.
}
