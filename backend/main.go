// main.go
package main

import (
	"GoAndDocker/backend/db"
	"GoAndDocker/backend/handlers"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	DB *db.Database
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if strings.HasPrefix(req.URL.Path, "/static/") {
		http.DefaultServeMux.ServeHTTP(w, req) // hand off to /static/ handler
		return
	}
	
	if req.URL.Path == "/" {
		http.ServeFile(w, req, "frontend/index.html")
		return
	}

	switch req.URL.Path {
	case "/api/register":
		handlers.RegisterHandler(w, req, r.DB) // Pass DB here
	case "/api/login":
		handlers.LoginHandler(w, req, r.DB) // Pass DB here
	case "/api/verify-email":
		handlers.VerifyEmailHandler(w, req, r.DB) // Pass DB here
	default:
		http.NotFound(w, req)
	}
}

func main() {
	database := db.NewDatabase()
	router := &Router{DB: database}

	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
