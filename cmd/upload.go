// Package cmd contains the commands for the script manager CLI tool.
package cmd

import (
	"fmt"
	"scriptmanager/internal"

	"github.com/spf13/cobra"
)

// Global variables for command-line flags.
var fileName string // The name of the script to upload.

// init initializes the upload command and its flags.
func init() {
	// Define flags for the upload command.
	uploadCmd.Flags().StringVarP(&fileName, fileNameParamName, fileNameShortParamName, "", "The name of the script to upload")
	// Add the upload command to the root command.
	rootCmd.AddCommand(uploadCmd)
}

// uploadCmd defines the "upload" command for the CLI tool.
// This command uploads a script to the server.
var uploadCmd = &cobra.Command{
	Use:   "upload",                        // Command name.
	Short: "Upload a script to the server", // Short description of the command.
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the script name from the command-line flags.
		scriptName, _ := cmd.Flags().GetString(fileNameParamName)

		// Load the script directory for the specified script.
		scriptDir, err := internal.LoadScriptDirectory(scriptName)
		if err != nil {
			return fmt.Errorf("could not load the script directory; %v", err)
		}
		defer scriptDir.ScriptReader.Close()

		// Prompt the user for credentials.
		username, password, err := internal.AskCredentials()
		if err != nil {
			return fmt.Errorf("could not read the username and password; %v", err)
		}

		// Log in to the server.
		c, err := internal.Login(scriptDir.ScriptInfo.ServerAddr, username, password)
		if err != nil {
			return fmt.Errorf("could not connect to the server; %v", err)
		}
		defer c.Quit()

		// Upload the file to the server.
		err = c.Stor(scriptDir.ScriptInfo.ServerPath, scriptDir.ScriptReader)
		if err != nil {
			return fmt.Errorf("could not upload the file; %v", err)
		}

		// Print a success message.
		fmt.Println("File uploaded successfully")

		return nil
	},
}
