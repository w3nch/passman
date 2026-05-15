package storage

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/w3nch/passman/internal/crypto"
	"github.com/w3nch/passman/internal/models"
)

const VaultFile = "data/vault.enc"
const ConfigFile = "data/config.json"

func LoadEntries(key []byte) ([]models.Entry, error) {
	var entries []models.Entry

	data, err := os.ReadFile(VaultFile)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return entries, nil
	}

	decrypted, err := crypto.Decrypt(data, key)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(decrypted, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func SaveEntries(entries []models.Entry, key []byte) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(data, key)
	if err != nil {
		return err
	}

	return os.WriteFile(VaultFile, encrypted, 0644)
}

func UsernameExists(entries []models.Entry, username string) bool {
	for _, entry := range entries {
		if strings.ToLower(entry.Username) == strings.ToLower(username) {
			return true
		}
	}
	return false
}

func GetAllEntries(key []byte) []models.Entry {
	entries, _ := LoadEntries(key)
	return entries
}

func DeleteEntry(service, username string, key []byte) bool {
	entries, err := LoadEntries(key)
	if err != nil {
		return false
	}

	for i, entry := range entries {
		if strings.ToLower(entry.Service) == strings.ToLower(service) &&
			strings.ToLower(entry.Username) == strings.ToLower(username) {
			entries = append(entries[:i], entries[i+1:]...)
			err = SaveEntries(entries, key)
			return err == nil
		}
	}

	return false
}

func UpdateEntry(service, username string, newUsername, newPassword string, key []byte) bool {
	entries, err := LoadEntries(key)
	if err != nil {
		return false
	}

	for i, entry := range entries {
		if strings.ToLower(entry.Service) == strings.ToLower(service) &&
			strings.ToLower(entry.Username) == strings.ToLower(username) {
			if newUsername != "" {
				entries[i].Username = newUsername
			}
			if newPassword != "" {
				entries[i].Password = newPassword
			}
			err = SaveEntries(entries, key)
			return err == nil
		}
	}

	return false
}

func SearchEntries(service, username string, key []byte) []models.Entry {
	entries, err := LoadEntries(key)
	if err != nil {
		return nil
	}

	var results []models.Entry

	for _, entry := range entries {
		serviceMatch := service == "" || strings.Contains(strings.ToLower(entry.Service), strings.ToLower(service))
		usernameMatch := username == "" || strings.Contains(strings.ToLower(entry.Username), strings.ToLower(username))

		if serviceMatch && usernameMatch {
			results = append(results, entry)
		}
	}

	return results
}

func LoadConfig() (models.VaultConfig, error) {
	var config models.VaultConfig

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func SaveConfig(config models.VaultConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0644)
}

func IsVaultSetup() bool {
	config, _ := LoadConfig()
	return config.IsSetup
}

func VaultExists() bool {
	_, err := os.Stat(VaultFile)
	return err == nil
}
