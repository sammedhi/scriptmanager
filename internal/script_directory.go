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

type ScriptDirectory struct {
	ScriptInfo   ScriptInfo
	ScriptReader io.Reader
}

func getAppCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	outputPath := filepath.Join(cacheDir, "scriptmanager")

	if err != nil {
		return "", err
	}

	return outputPath, nil
}

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

func OpenScript(scriptName string) error {
	cacheDir, err := getAppCacheDir()

	if err != nil {
		return fmt.Errorf("could not get the cache directory; %v", err)
	}

	scriptDirectoryPath := filepath.Join(cacheDir, scriptName)
	scriptInfoPath := filepath.Join(scriptDirectoryPath, "script_info.json")
	f, err := os.Open(scriptInfoPath)

	if err != nil {
		return fmt.Errorf("could not open the script info file; %v", err)
	}

	defer f.Close()

	var scriptInfo ScriptInfo
	bytes, _ := io.ReadAll(f)
	err = json.Unmarshal(bytes, &scriptInfo)

	if err != nil {
		return fmt.Errorf("could not decode the script info; %v", err)
	}

	scriptPath := filepath.Join(scriptDirectoryPath, scriptInfo.ScriptName+scriptInfo.ScriptExt)
	err = openFileWithDefaultProgram(scriptPath)

	if err != nil {
		return fmt.Errorf("could not open the file; %v", err)
	}

	return nil
}

func SaveScriptDirectory(scriptInfo ScriptInfo, script io.Reader) error {
	cacheDir, err := getAppCacheDir()
	outputPath := filepath.Join(cacheDir, scriptInfo.ScriptName)

	if err != nil {
		return fmt.Errorf("could not get the cache directory; %v", err)
	}

	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("directory creation failed; %v", err)
	}

	scriptDest := filepath.Join(outputPath, scriptInfo.ScriptName+scriptInfo.ScriptExt)
	scriptInfoDest := filepath.Join(outputPath, "script_info.json")

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

func LoadScriptDirectory(scriptName string) (*ScriptDirectory, error) {
	// Load the script info from the file
	// Load the file from the directory
	// Return a new ScriptDirectory object
	return nil, nil
}

func (sd *ScriptDirectory) Load() error {
	return nil
}

func (sd *ScriptDirectory) Save() error {
	return nil
}
