// Package internal contains internal utilities and logic for the script manager CLI tool.
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ScriptDirectory represents a directory containing a script and its metadata.
type ScriptDirectory struct {
	ScriptInfo   ScriptInfo    // Metadata about the script.
	ScriptReader io.ReadCloser // Reader for the script file.
}

// getAppCacheDir returns the cache directory path for the application.
// It creates a subdirectory named "scriptmanager" in the user's cache directory.
func getAppCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	outputPath := filepath.Join(cacheDir, "scriptmanager")

	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// getScriptDirectoryPath returns the path to the directory for a specific script.
// Parameters:
// - scriptName: The name of the script.
// Returns the full path to the script's directory.
func getScriptDirectoryPath(scriptName string) (string, error) {
	cacheDir, err := getAppCacheDir()

	if err != nil {
		return "", fmt.Errorf("could not get the cache directory; %v", err)
	}

	scriptDirectoryPath := filepath.Join(cacheDir, scriptName)

	return scriptDirectoryPath, nil
}

// openFileWithDefaultProgram opens a file using the default program for the current OS.
// Parameters:
// - filePath: The path to the file to open.
// Returns an error if the file cannot be opened.
func openFileWithDefaultProgram(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", filePath)
	}

	return cmd.Start()
}

// readScriptInfo reads the script metadata from the "script_info.json" file in the script's directory.
// Parameters:
// - scriptDirectoryPath: The path to the script's directory.
// Returns the ScriptInfo object or an error if the file cannot be read or parsed.
func readScriptInfo(scriptDirectoryPath string) (*ScriptInfo, error) {
	scriptInfoPath := filepath.Join(scriptDirectoryPath, "script_info.json")
	f, err := os.Open(scriptInfoPath)

	if err != nil {
		return nil, fmt.Errorf("could not open the script info file; %v", err)
	}

	defer f.Close()

	var scriptInfo ScriptInfo
	bytes, _ := io.ReadAll(f)
	err = json.Unmarshal(bytes, &scriptInfo)

	if err != nil {
		return nil, fmt.Errorf("could not decode the script info; %v", err)
	}

	return &scriptInfo, nil
}

// OpenScript opens a script file using the default program.
// Parameters:
// - scriptName: The name of the script to open.
// Returns an error if the script cannot be opened.
func OpenScript(scriptName string) error {
	scriptDirectoryPath, err := getScriptDirectoryPath(scriptName)

	if err != nil {
		return fmt.Errorf("could not get the script directory path; %v", err)
	}

	scriptInfo, err := readScriptInfo(scriptDirectoryPath)

	if err != nil {
		return fmt.Errorf("could not read the script info; `%`v", err)
	}

	scriptPath := filepath.Join(scriptDirectoryPath, scriptInfo.ScriptName+scriptInfo.ScriptExt)
	err = openFileWithDefaultProgram(scriptPath)

	if err != nil {
		return fmt.Errorf("could not open the file; %v", err)
	}

	return nil
}

// SaveScriptDirectory saves a script and its metadata to the local directory.
// Parameters:
// - scriptInfo: Metadata about the script.
// - script: The script content as an io.Reader.
// Returns an error if the script or metadata cannot be saved.
func SaveScriptDirectory(scriptInfo ScriptInfo, script io.Reader) error {
	scriptDirectoryPath, err := getScriptDirectoryPath(scriptInfo.ScriptName)

	if err != nil {
		return fmt.Errorf("could not get the script directory path; %v", err)
	}

	err = os.MkdirAll(scriptDirectoryPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("directory creation failed; %v", err)
	}

	scriptDest := filepath.Join(scriptDirectoryPath, scriptInfo.ScriptName+scriptInfo.ScriptExt)
	scriptInfoDest := filepath.Join(scriptDirectoryPath, "script_info.json")

	// Create the file
	dst, err := os.Create(scriptDest)
	if err != nil {
		return fmt.Errorf("file creation failed; %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, script)

	if err != nil {
		return fmt.Errorf("could not copy script to file; %v", err)
	}

	jsonInfo, err := json.Marshal(scriptInfo)
	if err != nil {
		return fmt.Errorf("could not marshal script info; %v", err)
	}

	err = os.WriteFile(scriptInfoDest, jsonInfo, 0644)
	if err != nil {
		return fmt.Errorf("could not write script info to file; %v", err)
	}

	return nil
}

// LoadScriptDirectory loads a script and its metadata from the local directory.
// Parameters:
// - scriptName: The name of the script to load.
// Returns a ScriptDirectory object or an error if the script cannot be loaded.
func LoadScriptDirectory(scriptName string) (*ScriptDirectory, error) {
	scriptDirectoryPath, err := getScriptDirectoryPath(scriptName)

	if err != nil {
		return nil, fmt.Errorf("could not get the script directory path; %v", err)
	}

	scriptInfo, err := readScriptInfo(scriptDirectoryPath)

	if err != nil {
		return nil, fmt.Errorf("could not read the script info; %v", err)
	}

	scriptPath := filepath.Join(scriptDirectoryPath, scriptInfo.ScriptName+scriptInfo.ScriptExt)
	scriptFile, err := os.Open(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("could not open the script file; %v", err)
	}

	return &ScriptDirectory{
		ScriptInfo:   *scriptInfo,
		ScriptReader: scriptFile,
	}, nil
}
