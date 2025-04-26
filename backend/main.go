// main.go
package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
	"os"
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

	// Update the query to insert username
	_, err = db.Exec("INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", u.Email, hashedPassword, u.Username)
	if err != nil {
		log.Println("Insert error:", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Database insert failed"}`, http.StatusInternalServerError)
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

	// Make sure the email credentials are set in environment variables
	if os.Getenv("EMAIL_USER") == "" || os.Getenv("EMAIL_PASSWORD") == "" {
		log.Fatal("Email credentials are not set in environment variables")
	}

	err := sendEmail("charles.ward.johnston@gmail.com", "Test Subject", "Test Body")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email sent successfully!")



	
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
