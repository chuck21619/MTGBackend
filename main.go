// main.go
package main

import (
	"GoAndDocker/backend/db"
	"GoAndDocker/backend/handlers"
	"GoAndDocker/backend/utils"
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

	// Serve frontend static files
	staticFS := http.FileServer(http.Dir("frontend/dist"))
	if _, err := fsOpen("frontend/dist" + req.URL.Path); err == nil {
		staticFS.ServeHTTP(w, req)
		return
	}

	// Fallback to index.html for React Router (SPA)
	http.ServeFile(w, req, "frontend/dist/index.html")
}

func fsOpen(path string) (http.File, error) {
	fs := http.Dir(".")
	return fs.Open(path)
}

func init() {
	//load local environment variables
	godotenv.Load()
}

func main() {
	utils.InitJWT()
	database := db.NewDatabase()
	router := &Router{DB: database}

	fs := http.FileServer(http.Dir("frontend/dist"))
	http.Handle("/", fs)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
