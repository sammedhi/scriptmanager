package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"scriptmanager/internal"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jlaffaye/ftp"
)

var serverAddr string
var filePath string
var username string = "ftp"
var password string = "!MonMotDePasse5"

func init() {
	fetchCmd.Flags().StringVar(&serverAddr, "addr", "", "The of ther server that hold the file")
	fetchCmd.Flags().StringVar(&filePath, "file_path", "", "The path where the file to fetch is located inside the server")
	rootCmd.AddCommand(fetchCmd)
}

func login(servAddr string, username string, password string) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(servAddr, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("could not connect to the server %s; %v", servAddr, err)
	}

	err = c.Login(username, password)

	if err != nil {
		return nil, fmt.Errorf("could not login to %s as %s; %v", servAddr, username, err)
	}

	return c, nil
}

func fetchFile(c *ftp.ServerConn, filePath string, fileDest string) error {
	file, err := c.Retr(filePath)

	if err != nil {
		return fmt.Errorf("could not retrieve the file '%s'; %v", filePath, err)
	}

	defer file.Close()

	// Ceate the directory if it does not exist
	dir := filepath.Dir(fileDest)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("directory creation failed; %v", err)
		}
	}

	// Create the file
	dst, err := os.Create(fileDest)
	if err != nil {
		return fmt.Errorf("file creation failed; %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		return fmt.Errorf("could not copy the file '%s' to '%s'; %v", filePath, fileDest, err)
	}

	return nil
}

func saveScriptInfo(scriptInfo internal.ScriptInfo, filePath string) error {
	jsonInfo, err := json.Marshal(scriptInfo)
	if err != nil {
		return fmt.Errorf("could not marshal script info; %v", err)
	}

	err = os.WriteFile(filePath, jsonInfo, 0644)
	if err != nil {
		return fmt.Errorf("could not write script info to file; %v", err)
	}

	return nil
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
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

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch the script at the target path",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the command line arguments
		serverAddr, _ := cmd.Flags().GetString("addr")
		filePath, err := cmd.Flags().GetString("file_path")

		if err != nil {
			return err
		}

		c, err := login(serverAddr, username, password)

		if err != nil {
			return err
		}

		defer c.Quit()

		// Try to read the file
		cacheDir, _ := os.UserCacheDir()
		outputPath := filepath.Join(cacheDir, "scriptmanager", fileNameWithoutExt(filePath), filepath.Base(filePath))
		err = fetchFile(c, filePath, outputPath)

		if err != nil {
			return fmt.Errorf("fetching the file failed; %v", err)
		}

		fmt.Printf("File '%s' fetched successfully to '%s'\n", filePath, outputPath)

		saveScriptInfo(internal.ScriptInfo{
			ScriptName: fileNameWithoutExt(filePath),
			ScriptExt:  filepath.Ext(filePath),
			ServerPath: filePath,
			ServerAddr: serverAddr,
		}, filepath.Join(cacheDir, "scriptmanager", fileNameWithoutExt(filePath), "script_info.json"))

		if err != nil {
			return fmt.Errorf("could not save script info; %v", err)
		}

		err = openFileWithDefaultProgram(outputPath)

		if err != nil {
			return fmt.Errorf("could not open the file; %v", err)
		}

		return nil
	},
}
