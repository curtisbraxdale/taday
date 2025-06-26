package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/curtisbraxdale/taday/internal/handlers"
	"github.com/curtisbraxdale/taday/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	apiCfg := handlers.ApiConfig{Queries: dbQueries, Platform: platform, Secret: secret}

	serveMux.HandleFunc("GET /api/ready", handlers.Ready)
	serveMux.HandleFunc("GET /api/users", apiCfg.GetUser)
	secure(serveMux, "GET /api/events", apiCfg.GetUserEvents, secret)
	serveMux.HandleFunc("POST /api/login", apiCfg.Login)

	server := http.Server{}
	server.Handler = serveMux
	server.Addr = ":8080"

	log.Fatal(server.ListenAndServe())
}

func secure(mux *http.ServeMux, methodAndPath string, handlerFunc http.HandlerFunc, secret string) {
	mux.Handle(methodAndPath, middleware.RequireAuth(secret, handlerFunc))
}
