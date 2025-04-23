package internal

import (
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"
)

func Login(servAddr string, username string, password string) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(servAddr, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("could not connect to the server %s; %v", servAddr, err)
	}

	err = c.Login(username, password)

	if err != nil {
		return nil, fmt.Errorf("could not login to %s as %s; %v", servAddr, username, err)
	}

	return c, nil
}
