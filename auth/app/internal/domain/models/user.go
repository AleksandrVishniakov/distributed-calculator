package models

type User struct {
	ID           uint64
	Login        string
	PasswordHash string
}
