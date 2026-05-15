package models

type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type VaultConfig struct {
	MasterHash string   `json:"master_hash"`
	Salt       string   `json:"salt"`
	Recovery   []string `json:"recovery_codes"`
	IsSetup    bool     `json:"is_setup"`
}
