package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/jlaffaye/ftp"
)

var serverAddr string
var filePath string

func init() {
	fetchCmd.Flags().StringVar(&serverAddr, "addr", "", "The of ther server that hold the file")
	fetchCmd.Flags().StringVar(&filePath, "file_path", "", "The path where the file to fetch is located inside the server")
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch the script at the target path",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverAddr, _ := cmd.Flags().GetString("addr")

		fmt.Printf("Try to reach the server %s\n", serverAddr)
		c, err := ftp.Dial(serverAddr)

		if err != nil {
			return err
		}

		filePath, err := cmd.Flags().GetString("file_path")
		if err != nil {
			fmt.Printf("There is no file_path argument\n")
			return err
		}

		err = c.Login("ftp", "!MonMotDePasse5")
		if err != nil {
			fmt.Println(err)
			return err

		}

		defer c.Quit()

		file, err := c.Retr(filePath)

		if err != nil {
			fmt.Printf("Could not retrieve the file %s\n", filePath)
			return err
		}

		defer file.Close()

		buf, err := io.ReadAll(file)

		if err != nil {
			fmt.Printf("Could not read the file %s\n", filePath)
			return err
		}

		fmt.Printf("file content : %s", buf)
		return nil
	},
}
