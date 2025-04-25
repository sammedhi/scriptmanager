// Package cmd contains the commands for the script manager CLI tool.
package cmd

import (
	"fmt"
	"path/filepath"
	"scriptmanager/internal"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

// Global variables for command-line flags.
var serverAddr string // The address of the server hosting the file.
var filePath string   // The path to the file on the server.

var filePathParamName = "file_path"     // Name of the file path parameter.
var filePathShortParamName = "p"        // Short name for the file path parameter.
var serverAddrParamName = "server_addr" // Name of the server address parameter.
var serverAddrShortParamName = "s"      // Short name for the server address parameter.

// init initializes the fetch command and its flags.
func init() {
	// Define flags for the fetch command.
	fetchCmd.Flags().StringVarP(&serverAddr, serverAddrParamName, serverAddrShortParamName, "", "The address of the server that holds the file")
	fetchCmd.Flags().StringVarP(&filePath, filePathParamName, filePathShortParamName, "", "The path where the file to fetch is located inside the server")
	// Add the fetch command to the root command.
	rootCmd.AddCommand(fetchCmd)
}

// fetchFile retrieves a file from the FTP server and saves it locally.
// Parameters:
// - c: The FTP server connection.
// - filePath: The path to the file on the server.
// Returns an error if the file cannot be retrieved or saved.
func fetchFile(c *ftp.ServerConn, filePath string) error {
	// Retrieve the file from the server.
	file, err := c.Retr(filePath)
	if err != nil {
		return fmt.Errorf("could not retrieve the file '%s'; %v", filePath, err)
	}
	defer file.Close()

	// Create a ScriptInfo object to store metadata about the script.
	scriptInfo := internal.ScriptInfo{
		ScriptName: fileNameWithoutExt(filePath),
		ScriptExt:  filepath.Ext(filePath),
		ServerPath: filePath,
		ServerAddr: serverAddr,
	}

	// Save the script to the local directory.
	err = internal.SaveScriptDirectory(scriptInfo, file)
	if err != nil {
		return fmt.Errorf("could not save the file '%s'; %v", filePath, err)
	}

	return nil
}

// fileNameWithoutExt extracts the file name without its extension.
// Parameters:
// - fileName: The full file name.
// Returns the file name without the extension.
func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

// fetchCmd defines the "fetch" command for the CLI tool.
// This command fetches a script from the server and saves it locally.
var fetchCmd = &cobra.Command{
	Use:   "fetch",                               // Command name.
	Short: "fetch the script at the target path", // Short description of the command.
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the server address and file path from the command-line flags.
		serverAddr, _ := cmd.Flags().GetString(serverAddrParamName)
		filePath, err := cmd.Flags().GetString(filePathParamName)
		if err != nil {
			return err
		}

		// Prompt the user for credentials.
		username, password, err := internal.AskCredentials()
		if err != nil {
			return fmt.Errorf("could not read the username and password; %v", err)
		}

		// Log in to the FTP server.
		c, err := internal.Login(serverAddr, username, password)
		if err != nil {
			return err
		}
		defer c.Quit()

		// Fetch the file from the server.
		err = fetchFile(c, filePath)
		if err != nil {
			return fmt.Errorf("fetching the file failed; %v", err)
		}

		// Print a success message.
		fmt.Printf("File '%s' fetched successfully\n", filePath)

		// Open the fetched script.
		err = internal.OpenScript(fileNameWithoutExt(filePath))
		if err != nil {
			return fmt.Errorf("could not open the file; %v", err)
		}

		return nil
	},
}
