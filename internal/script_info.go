package internal

type ScriptInfo struct {
	ScriptName string `json:"script_name"`
	ScriptExt  string `json:"script_ext"`
	ServerPath string `json:"server_path"`
	ServerAddr string `json:"server_addr"`
}
