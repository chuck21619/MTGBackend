package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"GoAndDocker/backend/models"
	"time"
	"os"
	"log"
)

var JwtKey []byte

func InitJWT() {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Fatal("JWT_SECRET_KEY environment variable not set")
	}
	JwtKey = []byte(key)
}

func GenerateAccessToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, expirationTime, err
}

func GenerateRefreshToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, expirationTime, err
}

func GenerateEmailVerificationToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}