package main

import (
	"github.com/chuck21619/MTGBackend/db"
	"github.com/chuck21619/MTGBackend/handlers"
	"github.com/chuck21619/MTGBackend/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Router struct {
	DB *db.Database
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api/") {
		switch req.URL.Path {
		case "/api/register":
			handlers.RegisterHandler(w, req, r.DB)
		case "/api/login":
			handlers.LoginHandler(w, req, r.DB)
		case "/api/logout":
			handlers.LogoutHandler(w, req, r.DB)
		case "/api/verify-email":
			handlers.VerifyEmailHandler(w, req, r.DB)
		case "/api/update-email":
			handlers.UpdateEmailHandler(w, req, r.DB)
		case "/api/refresh-token":
			handlers.RefreshTokenHandler(w, req, r.DB)
		case "/api/update-google-sheet":
			handlers.GoogleSheetHandler(w, req, r.DB)
		case "/api/populate":
			handlers.PopulateHandler(w, req, r.DB)
		case "/api/predict":
			handlers.PredictHandler(w, req, r.DB)
		default:
			http.NotFound(w, req)
		}
	}
}

func init() {
	loadEnv()
}

func loadEnv() {
	env := os.Getenv("ENV")
	if env != "docker" && env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}


func main() {
	utils.InitJWT()
	database := db.NewDatabase()
	router := &Router{DB: database}

	// CORS middleware setup
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "https://mtgfrontend.onrender.com", "http://frontend:5173"}, // Adjust this to your frontend URLs
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	log.Printf("Listening on :%s", port)
	// Wrap the router with the CORS handler
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, corsHandler.Handler(router)))
}
