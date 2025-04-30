package utils

import (
	"golang.org/x/crypto/bcrypt"
    "crypto/sha256"
    "encoding/hex"
)


func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashRefreshToken(refreshToken string) string {
    hash := sha256.Sum256([]byte(refreshToken))
    return hex.EncodeToString(hash[:])
}

func CheckRefreshTokenHash(token, hash string) bool {
	hashed := HashRefreshToken(token)
	return hashed == hash
}