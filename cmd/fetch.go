package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/jlaffaye/ftp"
)

var serverAddr string
var filePath string
var username string = "ftp"

func init() {
	fetchCmd.Flags().StringVar(&serverAddr, "addr", "", "The of ther server that hold the file")
	fetchCmd.Flags().StringVar(&filePath, "file_path", "", "The path where the file to fetch is located inside the server")
	rootCmd.AddCommand(fetchCmd)
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

		// Try to connect to the server
		c, err := ftp.Dial(serverAddr, ftp.DialWithTimeout(5*time.Second))

		if err != nil {
			return fmt.Errorf("could not connect to the server %s; %v", serverAddr, err)
		}

		// Try to login to the server
		err = c.Login(username, "!MonMotDePasse5")

		if err != nil {
			return fmt.Errorf("could not login to %s as %s; %v", serverAddr, username, err)
		}

		defer c.Quit()

		fmt.Printf("successfully logged as %s\n", username)

		// Try to retrieve the file
		fmt.Printf("try to retrieve the file '%s'\n", filePath)
		file, err := c.Retr(filePath)

		if err != nil {
			return fmt.Errorf("could not retrieve the file '%s'; %v", filePath, err)
		}

		defer file.Close()
		fmt.Printf("successfully retrieved \n")

		// Try to read the file
		fmt.Printf("try to read the file %s\n", filePath)
		buf, err := io.ReadAll(file)

		if err != nil {
			fmt.Printf("Could not read the file %s\n", filePath)
			return err
		}

		fmt.Printf("file content : %s", buf)
		return nil
	},
}
