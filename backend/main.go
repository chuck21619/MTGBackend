// main.go
package main

import (
	"log"
	"net/http"
    "GoAndDocker/backend/db"
	"GoAndDocker/backend/handlers"
)

type Router struct {
	DB *db.Database
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
    case "/register":
        handlers.RegisterHandler(w, req, r.DB)  // Pass DB here
    case "/login":
        handlers.LoginHandler(w, req, r.DB)  // Pass DB here
    case "/verify-email":
        handlers.VerifyEmailHandler(w, req, r.DB)  // Pass DB here
    default:
        http.NotFound(w, req)
    }
}

func main() {

	log.Println("IS IT GETTING STUCL")
    database := db.NewDatabase()
    router := &Router{DB: database}

	log.Println("Serving files from: ", http.Dir("frontend"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URL: ", r.URL)
		http.FileServer(http.Dir("frontend")).ServeHTTP(w, r)
	})
	// http.Handle("/", http.FileServer(http.Dir("frontend")))
    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}