package models

type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
}
