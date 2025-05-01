package handlers

import (
	"GoAndDocker/db"
	"GoAndDocker/models"
	"GoAndDocker/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Missing refresh token cookie")
		return
	}
	refreshToken := cookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	storedHash, err := database.GetRefreshTokenHash(claims.Username)
	if err != nil || !utils.CheckRefreshTokenHash(refreshToken, storedHash) {
		utils.WriteJSONMessage(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	newAccessToken, _, err := utils.GenerateAccessToken(claims.Username)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to create access token")
		return
	}

	newRefreshToken, refreshExpirationTime, err := utils.GenerateRefreshToken(claims.Username)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	hashedNewRefresh := utils.HashRefreshToken(newRefreshToken)
	if err := database.StoreRefreshToken(claims.Username, hashedNewRefresh); err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Failed to store refresh token")
		return
	}
	
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  refreshExpirationTime,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/refresh-token",
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
	})
}
