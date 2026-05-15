package input

import "os"

func UserInputs() (string, string, string, string) {

	command := os.Args[1]
	service := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	return command, service, username, password
}
