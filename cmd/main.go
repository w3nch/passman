package main

import (
	"fmt"
	"strings"

	"github.com/w3nch/passman/internal/generator"
	"github.com/w3nch/passman/internal/input"
	"github.com/w3nch/passman/internal/models"
	"github.com/w3nch/passman/internal/storage"
)

func main() {

	command, service, username, password := input.UserInputs()

	switch command {

	case "add":

		if strings.ToLower(password) == "random" {
			password = generator.Generate(16)
		}

		entry := models.Entry{
			Service:  service,
			Username: username,
			Password: password,
		}

		entries, err := storage.LoadEntries()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		exists := storage.UsernameExists(entries, username)

		if exists {
			fmt.Println("Username already exists")
			return
		}
		entries = append(entries, entry)

		err = storage.SaveEntries(entries)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Entry saved")

	case "list":
		fmt.Println("List command")
	case "search":
		fmt.Println("List command")
	case "edit":
		fmt.Println("List command")
	case "delete":
		fmt.Println("List command")
	default:
		fmt.Println("Unknown command")
	}
}
