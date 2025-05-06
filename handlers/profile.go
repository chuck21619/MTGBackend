package handlers

import (
	"net/http"
	"strings"
	"fmt"
	"encoding/json"

	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/utils"
	"github.com/chuck21619/MTGBackend/models"

	"github.com/golang-jwt/jwt/v5"
)

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
