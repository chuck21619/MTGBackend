package handlers

import (
	"encoding/json"
	"net/http"
	"GoAndDocker/db"
	"GoAndDocker/utils"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	if r.Method != http.MethodPost {
		utils.WriteJSONMessage(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := database.DeleteRefreshToken(req.Username)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Error deleting refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	utils.WriteJSONMessage(w, http.StatusOK, "Logged out successfully")
}