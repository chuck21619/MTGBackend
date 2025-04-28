package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"

	"GoAndDocker/backend/models"
	"GoAndDocker/backend/utils"
	"GoAndDocker/backend/db"
	"golang.org/x/crypto/bcrypt"
	//"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key_here")

func RegisterHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Generate the verification token
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		log.Println("Error generating token:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	// Generate token expiration time (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)
	// Insert the new user into the database with the verification token and emailVerified set to false
	_, err = database.Exec("INSERT INTO users (email, password, username, verification_token, email_verified, verification_token_expires_at) VALUES ($1, $2, $3, $4, $5, $6)", u.Email, hashedPassword, u.Username, token, false, expirationTime)
	if err != nil {
		log.Println("Insert error:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Database insert failed"}`, http.StatusInternalServerError)
		return
	}
	// Send the verification email
	err = utils.SendVerificationEmail(u.Email, token)
	if err != nil {
		log.Println("Failed to send verification email:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Registration successful, but failed to send email"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Please verify your email to log in"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	var storedHash string
	err = database.QueryRow("SELECT password, email_verified FROM users WHERE username = $1", u.Username).Scan(&storedHash, &u.Email_verified)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(u.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if !u.Email_verified {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Email has not been verified"}`, http.StatusNotFound)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// refresh token valid for 30 days
	refreshExpirationTime := time.Now().Add(30 * 24 * time.Hour)
	refreshClaims := &models.Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Set refresh token as secure HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  refreshExpirationTime,
		HttpOnly: true,
		Secure:   true, // true if using HTTPS
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// Send access token in response body
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenString,
		"message": "Login successful",
	})
}