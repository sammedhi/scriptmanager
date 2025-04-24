package cmd

import (
	"fmt"
	"scriptmanager/internal"

	"github.com/spf13/cobra"
)

var fileNameParamName = "file_name"
var fileNameShortParamName = "n"

func init() {
	openCmd.Flags().StringVarP(&fileName, fileNameParamName, fileNameShortParamName, "", "The name of the file to upload")
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open one of the script fetched previously",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptName, _ := cmd.Flags().GetString(fileNameParamName)

		scriptDir, err := internal.LoadScriptDirectory(scriptName)

		if err != nil {
			return fmt.Errorf("could not load the script directory; %v", err)
		}

		defer scriptDir.ScriptReader.Close()

		err = internal.OpenScript(scriptDir.ScriptInfo.ServerPath)

		if err != nil {
			return fmt.Errorf("could not open the file; %v", err)
		}

		return nil
	},
}
