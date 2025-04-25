// Package internal contains internal utilities and logic for the script manager CLI tool.
package internal

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// AskCredentials prompts the user for a username and password.
// Returns the username, password, or an error if input fails.
func AskCredentials() (string, string, error) {
	var username string

	// Prompt for the username.
	fmt.Print("username: ")
	fmt.Scanln(&username)

	// Prompt for the password.
	fmt.Print("password: ")
	pwd, err := term.ReadPassword(int(syscall.Stdin))
	password := string(pwd)

	if err != nil {
		fmt.Println("Error reading password:", err)
		return "", "", err
	}

	return username, password, nil
}
