package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/models"
	"github.com/chuck21619/MTGBackend/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Game map[string]string

func PopulateHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
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

	google_sheet, err := database.GetGoogleSheet(claims.Username)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	println("google sheet url: ", google_sheet)

	players, decks, err := getUniquePlayersAndDecks(google_sheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := map[string][]string{
		"players": players,
		"decks":   decks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func PredictHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
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

	type Selection struct {
		Player string `json:"player"`
		Deck   string `json:"deck"`
	}
	type PredictRequest struct {
		Selections []Selection `json:"selections"`
	}

	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Received selections:")
	for i, s := range req.Selections {
		fmt.Printf("Player %d: %s, Deck: %s\n", i+1, s.Player, s.Deck)
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Internal error")
	}

	microserviceURL := os.Getenv("MICROSERVICE_URL") + "/predict"
	// Call the microservice
	resp, err := http.Post(microserviceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "Failed to contact microservice", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Return the raw response from the microservice to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func TrainHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
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

	microserviceURL := os.Getenv("MICROSERVICE_URL") + "/train"
	google_sheet, err := database.GetGoogleSheet(claims.Username)
	if err != nil {
		utils.WriteJSONMessage(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	
	jsonBody := []byte(fmt.Sprintf(`{"url": "%s"}`, google_sheet))

	resp, err := http.Post(microserviceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "Failed to contact microservice", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func fetchCSVData(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	var data [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, record)
	}
	return data, nil
}

func getUniquePlayersAndDecks(url string) ([]string, []string, error) {
	data, err := fetchCSVData(url)
	if err != nil {
		return nil, nil, err
	}

	if len(data) < 2 {
		return nil, nil, fmt.Errorf("not enough data")
	}

	headers := data[0]
	var playerNames []string
	deckSet := make(map[string]struct{})

	for _, header := range headers {
		if header != "winner" {
			playerNames = append(playerNames, header)
		}
	}

	for _, row := range data[1:] {
		for i, val := range row {
			if headers[i] != "winner" && val != "" {
				deckSet[val] = struct{}{}
			}
		}
	}

	var deckNames []string
	for deck := range deckSet {
		deckNames = append(deckNames, deck)
	}

	sort.Strings(playerNames)
	sort.Strings(deckNames)

	return playerNames, deckNames, nil
}
