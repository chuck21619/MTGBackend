package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"GoAndDocker/backend/models"
	"GoAndDocker/backend/utils"
	"GoAndDocker/backend/db"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	if r.Method != http.MethodPost {
		utils.WriteJSONMessage(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	_, err = database.Exec("INSERT INTO users (email, password, username, verification_token, email_verified, verification_token_expires_at) VALUES ($1, $2, $3, $4, $5, $6)", u.Email, hashedPassword, u.Username, token, false, expirationTime)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Database insert failed")
		return
	}

	err = utils.SendVerificationEmail(u.Email, token)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Registration successful, but failed to send email")
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "Please verify your email to log in")
}

