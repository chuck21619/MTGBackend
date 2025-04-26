// main.go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func sendEmail(to string, subject string, body string) error {

	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func generateVerificationToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

func sendVerificationEmail(userEmail, token string) error {
	verificationURL := fmt.Sprintf("https://goanddocker.onrender.com/verify-email?token=%s", token)
	subject := "Please verify your email address"
	body := fmt.Sprintf("Click this link to verify your email address: %s", verificationURL)
	return sendEmail(userEmail, subject, body)
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	var userID int
	var emailVerified bool
	var tokenExpiration time.Time
	err := db.QueryRow("SELECT id, email_verified, token_expiration FROM users WHERE verification_token = $1", token).Scan(&userID, &emailVerified, &tokenExpiration)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusNotFound)
		return
	}

	if emailVerified {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusNotFound)
		return
	}

	_, err = db.Exec("UPDATE users SET email_verified = TRUE WHERE id = $1", userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Failed to verify email"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Your email has been successfully verified. You can now log in."}`))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var u User
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
	token, err := generateVerificationToken()
	if err != nil {
		log.Println("Error generating token:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	// Generate token expiration time (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)
	// Insert the new user into the database with the verification token and emailVerified set to false
	_, err = db.Exec("INSERT INTO users (email, password, username, verification_token, email_verified, token_expires_at) VALUES ($1, $2, $3, $4, $5, $6)", u.Email, hashedPassword, u.Username, token, false, expirationTime)
	if err != nil {
		log.Println("Insert error:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Database insert failed"}`, http.StatusInternalServerError)
		return
	}
	// Send the verification email
	err = sendVerificationEmail(u.Email, token)
	if err != nil {
		log.Println("Failed to send verification email:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Registration successful, but failed to send email"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	var storedHash string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", u.Username).Scan(&storedHash)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func main() {
	initDB()

	http.Handle("/", http.FileServer(http.Dir("frontend")))

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/verify-email", verifyEmailHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
