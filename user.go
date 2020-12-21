package main

import "time"

// User ...
type User struct {
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	CreatedAt  time.Time `json:"createdAt"`
	Salt       string    `json:"salt"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	AuthType   int       `json:"authType"`
	Token      string    `json:"token"`
	ResetToken string    `json:"resetToken"`
}
