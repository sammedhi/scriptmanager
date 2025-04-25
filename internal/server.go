// Package internal contains internal utilities and logic for the script manager CLI tool.
package internal

import (
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"
)

// Login connects to an FTP server and logs in using the provided credentials.
// Parameters:
// - servAddr: The address of the FTP server.
// - username: The username for authentication.
// - password: The password for authentication.
// Returns an FTP server connection or an error if the connection or login fails.
func Login(servAddr string, username string, password string) (*ftp.ServerConn, error) {
	// Connect to the FTP server with a timeout.
	c, err := ftp.Dial(servAddr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("could not connect to the server %s; %v", servAddr, err)
	}

	// Log in to the FTP server.
	err = c.Login(username, password)
	if err != nil {
		return nil, fmt.Errorf("could not login to %s as %s; %v", servAddr, username, err)
	}

	return c, nil
}
