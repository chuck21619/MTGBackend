// main.go
package main

import (
	"GoAndDocker/db"
	"GoAndDocker/handlers"
	"GoAndDocker/utils"
	"log"
	"net/http"
	"strings"
	"github.com/joho/godotenv"
)

type Router struct {
	DB *db.Database
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle API routes
	if strings.HasPrefix(req.URL.Path, "/api/") {
		switch req.URL.Path {
		case "/api/register":
			handlers.RegisterHandler(w, req, r.DB)
		case "/api/login":
			handlers.LoginHandler(w, req, r.DB)
		case "/api/verify-email":
			handlers.VerifyEmailHandler(w, req, r.DB)
		case "/api/update-email":
			handlers.UpdateEmailHandler(w, req, r.DB)
		case "/api/refresh-token":
			handlers.RefreshTokenHandler(w, req, r.DB)
		default:
			http.NotFound(w, req)
		}
		return
	}
}

func init() {
	//load local environment variables
	godotenv.Load()
}

func main() {
	utils.InitJWT()
	database := db.NewDatabase()
	router := &Router{DB: database}

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
