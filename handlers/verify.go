package handlers

import (
	"log"
	"net/http"
	"github.com/chuck21619/MTGBackend/db"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	var userID int
	var emailVerified bool
	err := database.QueryRow("SELECT id, email_verified FROM users WHERE verification_token = $1", token).Scan(&userID, &emailVerified)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusNotFound)
		return
	}

	if emailVerified {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Email is already verified"}`, http.StatusNotFound)
		return
	}

	_, err = database.Exec("UPDATE users SET email_verified = TRUE WHERE id = $1", userID)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Failed to verify email"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Your email has been successfully verified. You can now log in."}`))
}
