package cmd

import (
	"fmt"
	"scriptmanager/internal"

	"github.com/spf13/cobra"
)

var scriptName string

func init() {
	uploadCmd.Flags().StringVar(&scriptName, "script_name", "", "The name of the script to upload")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a script to the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptName, _ := cmd.Flags().GetString("script_name")

		scriptDir, err := internal.LoadScriptDirectory(scriptName)

		if err != nil {
			return fmt.Errorf("could not load the script directory; %v", err)
		}

		defer scriptDir.ScriptReader.Close()

		// Connect to the server
		c, err := internal.Login(scriptDir.ScriptInfo.ServerAddr, username, password)

		if err != nil {
			return fmt.Errorf("could not connect to the server; %v", err)
		}

		// Upload the file
		err = c.Stor(scriptDir.ScriptInfo.ServerPath, scriptDir.ScriptReader)

		if err != nil {
			return fmt.Errorf("could not upload the file; %v", err)
		}

		return nil
	},
}
