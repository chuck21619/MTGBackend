// main.go
package main

import (
	"log"
	"net/http"
    "GoAndDocker/backend/db"
	"GoAndDocker/backend/handlers"
	"os"
)

type Router struct {
	DB *db.Database
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	
	staticFilePath := "frontend" + req.URL.Path
    if fileExists(staticFilePath) {
        log.Println("Serving static file: ", staticFilePath)
        http.FileServer(http.Dir("frontend")).ServeHTTP(w, req)
        return
    }
	 
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
    database := db.NewDatabase()
    router := &Router{DB: database}
    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func fileExists(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !info.IsDir()
}