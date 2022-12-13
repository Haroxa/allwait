package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Wait struct {
	gorm.Model
	Title     string
	Number    string
	Wtime     int
	Address   string
	Transport string
}
