package main

import (
	"fmt"
	"github.com/w3nch/passman.git/internal/generator"
	"os"
)

func main() {
	command := os.Args[1]
	service := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]
	fmt.Println(command, service, username, password)
	fmt.Println(password)
	if password == "random" {
		fmt.Println(generator.Generate(16))
	}
}
