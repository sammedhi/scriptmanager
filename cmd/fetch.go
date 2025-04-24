package cmd

import (
	"fmt"
	"path/filepath"
	"scriptmanager/internal"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jlaffaye/ftp"
)

var serverAddr string
var filePath string
var filePathParamName = "file_path"
var filePathShortParamName = "p"

var serverAddrParamName = "server_addr"
var serverAddrShortParamName = "s"

func init() {
	fetchCmd.Flags().StringVarP(&serverAddr, serverAddrParamName, serverAddrShortParamName, "", "The address of the server that hold the file")
	fetchCmd.Flags().StringVarP(&filePath, filePathParamName, filePathShortParamName, "", "The path where the file to fetch is located inside the server")
	rootCmd.AddCommand(fetchCmd)
}

func fetchFile(c *ftp.ServerConn, filePath string) error {
	file, err := c.Retr(filePath)

	if err != nil {
		return fmt.Errorf("could not retrieve the file '%s'; %v", filePath, err)
	}

	defer file.Close()

	scriptInfo := internal.ScriptInfo{
		ScriptName: fileNameWithoutExt(filePath),
		ScriptExt:  filepath.Ext(filePath),
		ServerPath: filePath,
		ServerAddr: serverAddr,
	}

	err = internal.SaveScriptDirectory(scriptInfo, file)

	if err != nil {
		return fmt.Errorf("could not save the file '%s'; %v", filePath, err)
	}

	return nil
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch the script at the target path",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the command line arguments
		serverAddr, _ := cmd.Flags().GetString(serverAddrParamName)
		filePath, err := cmd.Flags().GetString(filePathParamName)

		if err != nil {
			return err
		}

		username, password, err := internal.AskCredentials()

		if err != nil {
			return fmt.Errorf("could not read the username and password; %v", err)
		}

		c, err := internal.Login(serverAddr, username, password)

		if err != nil {
			return err
		}

		defer c.Quit()

		err = fetchFile(c, filePath)

		if err != nil {
			return fmt.Errorf("fetching the file failed; %v", err)
		}

		fmt.Printf("File '%s' fetched successfully\n", filePath)

		if err != nil {
			return fmt.Errorf("could not save script info; %v", err)
		}

		err = internal.OpenScript(fileNameWithoutExt(filePath))

		if err != nil {
			return fmt.Errorf("could not open the file; %v", err)
		}

		return nil
	},
}
