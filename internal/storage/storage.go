package storage

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/w3nch/passman/internal/models"
)

const VaultFile = "data/vault.json"

func LoadEntries() ([]models.Entry, error) {

	var entries []models.Entry

	data, err := os.ReadFile(VaultFile)

	if err != nil {

		if os.IsNotExist(err) {
			return entries, nil
		}

		return nil, err
	}

	err = json.Unmarshal(data, &entries)

	if err != nil {
		return nil, err
	}

	return entries, nil
}

func SaveEntries(entries []models.Entry) error {

	data, err := json.MarshalIndent(entries, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(VaultFile, data, 0644)
}

func UsernameExists(entries []models.Entry, username string) bool {

	for _, entry := range entries {

		if strings.ToLower(entry.Username) == strings.ToLower(username) {
			return true
		}
	}

	return false
}
