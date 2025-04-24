package internal

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func AskCredentials() (string, string, error) {
	var username string
	fmt.Print("username: ")
	fmt.Scanln(&username)

	fmt.Print("password: ")
	pwd, err := term.ReadPassword(int(syscall.Stdin))
	password := string(pwd)

	if err != nil {
		fmt.Println("Error reading password:", err)
		return "", "", err
	}

	return username, password, nil
}
