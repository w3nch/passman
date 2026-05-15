package input

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func UserInputs() (string, string, string, string) {

	command := os.Args[1]
	service := ""
	username := ""
	password := ""

	if len(os.Args) > 2 {
		service = os.Args[2]
	}
	if len(os.Args) > 3 {
		username = os.Args[3]
	}
	if len(os.Args) > 4 {
		password = os.Args[4]
	}

	return command, service, username, password
}

func ReadPassword(prompt string) string {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(password)
}

func ReadInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
