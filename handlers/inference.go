package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"encoding/json"

	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/models"
	"github.com/chuck21619/MTGBackend/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Game map[string]string

func PopulateHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {

	print("LOOK AT ME IM MR MEESEIX")

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
	fmt.Println(players)
	fmt.Println(decks)

	response := map[string][]string{
		"players": players,
		"decks":   decks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

func generateDataset(url string) ([]Game, error) {
	data, err := fetchCSVData(url)
	if err != nil {
		return nil, err
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("not enough rows")
	}

	headers := data[0]
	var games []Game

	for _, row := range data[1:] {
		game := Game{}
		for i, val := range row {
			header := headers[i]
			if header != "winner" && val != "" {
				game[header] = val
			}
		}
		// Add winner
		for i, val := range row {
			if headers[i] == "winner" {
				game["winner"] = val
				break
			}
		}
		games = append(games, game)
	}
	return games, nil
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
