package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/utils"
)

func ProfileInfo(w http.ResponseWriter, r *http.Request, database *db.Database) {
	claims, ok := utils.ValidateJWT(w, r)
	if !ok {
		return
	}

	user, err := database.GetProfileInfo(claims.Username)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"email":          user.Email,
		"googleSheetUrl": user.GoogleSheet,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateEmailHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	claims, ok := utils.ValidateJWT(w, r)
	if !ok {
		return
	}

	var body struct {
		NewEmail string `json:"new_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.NewEmail == "" {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := database.UpdateUserEmail(claims.Username, body.NewEmail)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to update email")
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "Email updated successfully")
}

func GoogleSheetHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	claims, ok := utils.ValidateJWT(w, r)
	if !ok {
		return
	}

	var body struct {
		NewGoogleSheet string `json:"new_google_sheet"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.NewGoogleSheet == "" {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := database.UpdateGoogleSheet(claims.Username, body.NewGoogleSheet)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to update google sheet")
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "Google sheet updated successfully")
}
