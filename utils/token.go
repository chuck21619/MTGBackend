package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/chuck21619/MTGBackend/models"
	"time"
	"os"
	"log"
	"fmt"
	"net/http"
	"strings"
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


func ValidateJWT(w http.ResponseWriter, r *http.Request) (*models.Claims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		WriteJSONMessage(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return nil, false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		WriteJSONMessage(w, http.StatusUnauthorized, "Invalid or expired token")
		return nil, false
	}

	return claims, true
}