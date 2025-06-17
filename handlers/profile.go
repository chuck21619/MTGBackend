package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/models"
	"github.com/chuck21619/MTGBackend/utils"

	"github.com/golang-jwt/jwt/v5"
)

func ProfileInfo(w http.ResponseWriter, r *http.Request, database *db.Database) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Invalid or expired token")
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	var body struct {
		NewEmail string `json:"new_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.NewEmail == "" {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = database.UpdateUserEmail(claims.Username, body.NewEmail)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to update email")
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "Email updated successfully")
}

func GoogleSheetHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	var body struct {
		NewGoogleSheet string `json:"new_google_sheet"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.NewGoogleSheet == "" {
		utils.WriteJSONMessage(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = database.UpdateGoogleSheet(claims.Username, body.NewGoogleSheet)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to update google sheet")
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "Google sheet updated successfully")
}
