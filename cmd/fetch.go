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
var username string = "ftpuser"
var password string = "gigigi"

func init() {
	fetchCmd.Flags().StringVar(&serverAddr, "addr", "", "The of ther server that hold the file")
	fetchCmd.Flags().StringVar(&filePath, "file_path", "", "The path where the file to fetch is located inside the server")
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
		serverAddr, _ := cmd.Flags().GetString("addr")
		filePath, err := cmd.Flags().GetString("file_path")

		if err != nil {
			return err
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
