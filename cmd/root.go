package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "scm",
	Short:        "Script manager is a simple tool to manage script from a remote server",
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}
