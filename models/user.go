package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Email                         string    `json:"email"`
	Password                      string    `json:"password"`
	Username                      string    `json:"username"`
	Verification_token            string    `json:"verification_token"`
	Email_verified                bool      `json:"email_verified"`
	Verification_token_expires_at time.Time `json:"verification_token_expires_at"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
