// Package cmd contains the commands for the script manager CLI tool.
package cmd

import (
	"fmt"
	"scriptmanager/internal"

	"github.com/spf13/cobra"
)

// Global variables for command-line flags.
var fileNameParamName = "file_name" // Name of the file name parameter.
var fileNameShortParamName = "n"    // Short name for the file name parameter.

// init initializes the open command and its flags.
func init() {
	// Define flags for the open command.
	openCmd.Flags().StringVarP(&fileName, fileNameParamName, fileNameShortParamName, "", "The name of the file to open")
	// Add the open command to the root command.
	rootCmd.AddCommand(openCmd)
}

// openCmd defines the "open" command for the CLI tool.
// This command opens a script that was previously fetched.
var openCmd = &cobra.Command{
	Use:   "open",                                       // Command name.
	Short: "open one of the scripts fetched previously", // Short description of the command.
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the script name from the command-line flags.
		scriptName, _ := cmd.Flags().GetString(fileNameParamName)

		// Load the script directory for the specified script.
		scriptDir, err := internal.LoadScriptDirectory(scriptName)
		if err != nil {
			return fmt.Errorf("could not load the script directory; %v", err)
		}
		defer scriptDir.ScriptReader.Close()

		// Open the script using the internal package.
		err = internal.OpenScript(scriptDir.ScriptInfo.ServerPath)
		if err != nil {
			return fmt.Errorf("could not open the file; %v", err)
		}

		return nil
	},
}
