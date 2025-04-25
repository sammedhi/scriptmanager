// Package cmd contains the commands for the script manager CLI tool.
package cmd

import "github.com/spf13/cobra"

// rootCmd is the root command for the CLI tool.
// It serves as the entry point for all subcommands.
var rootCmd = &cobra.Command{
	Use:          "scm",                                                                    // Command name.
	Short:        "Script manager is a simple tool to manage scripts from a remote server", // Short description of the tool.
	SilenceUsage: true,                                                                     // Prevents usage message from being printed on errors.
}

// Execute runs the root command.
// This function is called from the main package to start the CLI tool.
func Execute() error {
	return rootCmd.Execute()
}
