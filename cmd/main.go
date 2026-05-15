package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/w3nch/passman/internal/crypto"
	"github.com/w3nch/passman/internal/generator"
	"github.com/w3nch/passman/internal/input"
	"github.com/w3nch/passman/internal/models"
	"github.com/w3nch/passman/internal/storage"
)

var isUnlocked = false
var isSetup = false
var encryptionKey []byte

func main() {
	isSetup = storage.IsVaultSetup()

	if len(os.Args) > 1 && os.Args[1] == "setup" {
		if isSetup {
			fmt.Println("Vault is already set up. Use 'unlock' to access.")
			return
		}
		setupVault()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "unlock" {
		if unlockVault() {
			runMainMenu()
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "reset" {
		resetVault()
		return
	}

	if !isSetup {
		fmt.Println("Vault not set up. Run: passman setup")
		return
	}

	if !isUnlocked {
		fmt.Println("Vault is locked. Run: passman unlock")
		return
	}

	runMainMenu()
}

func setupVault() {
	fmt.Println("=== First Time Setup ===")
	fmt.Println("Creating your master password...")
	fmt.Println()

	password := input.ReadPassword("Enter master password: ")
	if password == "" {
		fmt.Println("Password cannot be empty")
		return
	}

	confirm := input.ReadPassword("Confirm master password: ")
	if password != confirm {
		fmt.Println("Passwords do not match")
		return
	}

	salt := crypto.GenerateSalt()
	hash := crypto.HashPassword(password, salt)
	encryptionKey = crypto.DeriveKey(password, salt)
	recoveryCodes := crypto.GenerateRecoveryCodes(5)

	fmt.Println()
	fmt.Println("=== RECOVERY CODES - SAVE THESE! ===")
	for i, code := range recoveryCodes {
		fmt.Printf("%d. %s\n", i+1, code)
	}
	fmt.Println("====================================")
	fmt.Println("You will need these to reset your password if forgotten.")
	fmt.Println()

	config := models.VaultConfig{
		MasterHash: hash,
		Salt:       salt,
		Recovery:   recoveryCodes,
		IsSetup:    true,
	}

	err := storage.SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	entries := []models.Entry{}
	err = storage.SaveEntries(entries, encryptionKey)
	if err != nil {
		fmt.Println("Error initializing vault:", err)
		return
	}

	fmt.Println("Vault setup complete!")
	isUnlocked = true
	runMainMenu()
}

func unlockVault() bool {
	password := input.ReadPassword("Enter master password: ")

	config, err := storage.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return false
	}

	if crypto.VerifyPassword(password, config.Salt, config.MasterHash) {
		encryptionKey = crypto.DeriveKey(password, config.Salt)
		isUnlocked = true
		fmt.Println("Vault unlocked!")
		return true
	}

	fmt.Println("Incorrect password")
	return false
}

func resetVault() {
	fmt.Println("=== Password Reset ===")
	fmt.Println("Enter one of your recovery codes:")

	config, err := storage.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	code := input.ReadInput("Recovery code: ")
	code = strings.TrimSpace(code)

	valid := false
	for _, c := range config.Recovery {
		if c == code {
			valid = true
			break
		}
	}

	if !valid {
		fmt.Println("Invalid recovery code")
		return
	}

	fmt.Println("Recovery code accepted!")
	fmt.Println("Enter your new master password:")

	password := input.ReadPassword("New password: ")
	if password == "" {
		fmt.Println("Password cannot be empty")
		return
	}

	confirm := input.ReadPassword("Confirm new password: ")
	if password != confirm {
		fmt.Println("Passwords do not match")
		return
	}

	salt := crypto.GenerateSalt()
	hash := crypto.HashPassword(password, salt)
	newRecoveryCodes := crypto.GenerateRecoveryCodes(5)

	config.MasterHash = hash
	config.Salt = salt
	config.Recovery = newRecoveryCodes

	err = storage.SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	fmt.Println()
	fmt.Println("=== NEW RECOVERY CODES ===")
	for i, code := range newRecoveryCodes {
		fmt.Printf("%d. %s\n", i+1, code)
	}
	fmt.Println("=========================")
	fmt.Println("Password reset complete!")
}

func runMainMenu() {
	fmt.Println("\nVault unlocked. Commands: add, list, search, edit, delete, lock, exit")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		service := ""
		username := ""
		password := ""

		if len(parts) > 1 {
			service = parts[1]
		}
		if len(parts) > 2 {
			username = parts[2]
		}
		if len(parts) > 3 {
			password = parts[3]
		}

		switch command {
		case "add":
			if service == "" || username == "" || password == "" {
				fmt.Println("Usage: add <service> <username> <password>")
				fmt.Println("Use 'random' for auto-generate (e.g., add github alice random)")
				continue
			}
			genUsername := strings.ToLower(username) == "random"
			genPassword := strings.ToLower(password) == "random"

			if genUsername && genPassword {
				username, password = generator.GenerateProfile()
			} else if genUsername {
				username, password = generator.GenerateProfile()
			} else if genPassword {
				_, password = generator.GenerateProfile()
			}

			if genUsername || genPassword {
				fmt.Printf("username: %s\n", username)
				fmt.Printf("password: %s\n", password)
			}

			entry := models.Entry{
				Service:  service,
				Username: username,
				Password: password,
			}

			entries, err := storage.LoadEntries(encryptionKey)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			exists := storage.UsernameExists(entries, username)
			if exists {
				fmt.Println("Username already exists")
				continue
			}
			entries = append(entries, entry)

			err = storage.SaveEntries(entries, encryptionKey)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Println("Entry saved")

		case "list":
			entries := storage.GetAllEntries(encryptionKey)
			if len(entries) == 0 {
				fmt.Println("No entries found")
				continue
			}
			fmt.Println("+----------+-------------------------+")
			fmt.Println("| Service  | Username                |")
			fmt.Println("+----------+-------------------------+")
			for _, entry := range entries {
				fmt.Printf("| %-8s | %-23s |\n", entry.Service, entry.Username)
			}
			fmt.Println("+----------+-------------------------+")

		case "search":
			if service == "" && username == "" {
				fmt.Println("Usage: search [service] [username]")
				fmt.Println("  search github     # by service")
				fmt.Println("  search alice       # by username")
				fmt.Println("  search github alice  # both")
				continue
			}
			results := storage.SearchEntries(service, username, encryptionKey)
			if len(results) == 0 {
				fmt.Println("No entries found")
				continue
			}
			fmt.Println("+----------+-------------------------+")
			fmt.Println("| Service  | Username                |")
			fmt.Println("+----------+-------------------------+")
			for _, entry := range results {
				fmt.Printf("| %-8s | %-23s |\n", entry.Service, entry.Username)
			}
			fmt.Println("+----------+-------------------------+")

		case "edit":
			if service == "" || username == "" {
				fmt.Println("Usage: edit <service> <username> [new_username] [new_password]")
				continue
			}
			newUsername := ""
			newPassword := ""
			if len(parts) > 3 {
				newUsername = parts[3]
			}
			if len(parts) > 4 {
				newPassword = parts[4]
			}
			if storage.UpdateEntry(service, username, newUsername, newPassword, encryptionKey) {
				fmt.Println("Entry updated")
			} else {
				fmt.Println("Entry not found")
			}

		case "delete":
			if service == "" || username == "" {
				fmt.Println("Usage: delete <service> <username>")
				continue
			}
			if storage.DeleteEntry(service, username, encryptionKey) {
				fmt.Println("Entry deleted")
			} else {
				fmt.Println("Entry not found")
			}

		case "lock":
			isUnlocked = false
			fmt.Println("Vault locked")
			return

		case "exit", "quit":
			fmt.Println("Goodbye!")
			isUnlocked = false
			return

		default:
			fmt.Println("Unknown command. Use: add, list, search, edit, delete, lock, exit")
		}
	}
}